// Copyright (c) 2015 Uber Technologies, Inc.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package thrift

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/apache/thrift/lib/go/thrift"
	tchannel "github.com/uber/tchannel-go"
	"golang.org/x/net/context"
)

type handler struct {
	standard       TChanServer
	streaming      TChanStreamingServer
	postResponseCB PostResponseCB
}

// Server handles incoming TChannel calls and forwards them to the matching TChanServer.
type Server struct {
	tcommon

	Channel       tchannel.Registrar
	log           tchannel.Logger
	mut           sync.RWMutex
	handlers      map[string]handler
	healthHandler *healthHandler
}

// NewServer returns a server that can serve thrift services over TChannel.
func NewServer(registrar tchannel.Registrar) *Server {
	healthHandler := newHealthHandler()
	server := &Server{
		Channel:       registrar,
		log:           registrar.Logger(),
		handlers:      make(map[string]handler),
		healthHandler: healthHandler,
	}

	server.Register(newTChanMetaServer(healthHandler))
	return server
}

func (s *Server) registerHandler(ts TChanServer, h handler, opts ...RegisterOption) {
	serviceName := ts.Service()
	for _, opt := range opts {
		opt.Apply(&h)
	}

	s.mut.Lock()
	s.handlers[serviceName] = h
	s.mut.Unlock()

	for _, m := range ts.Methods() {
		s.Channel.Register(s, serviceName+"::"+m)
	}

	if sts, ok := ts.(TChanStreamingServer); ok {
		if len(sts.StreamingMethods()) > 0 && h.streaming == nil {
			panic("Using Register when you should use RegisterStreaming")
		}
		for _, m := range sts.StreamingMethods() {
			s.Channel.Register(tchannel.HandlerFunc(s.HandleStreaming), serviceName+"::"+m)
		}
	}
}

// Register registers the given TChanServer to be called on any incoming call for its' services.
func (s *Server) Register(svr TChanServer, opts ...RegisterOption) {
	s.registerHandler(svr, handler{standard: svr}, opts...)
}

// RegisterStreaming registers the given TChanStreamingServer to be called on any incoming call
// for its' services.
func (s *Server) RegisterStreaming(svr TChanStreamingServer, opts ...RegisterOption) {
	svr.SetCommon(s.tcommon)
	s.registerHandler(svr, handler{standard: svr, streaming: svr}, opts...)
}

// RegisterHealthHandler uses the user-specified function f for the Health endpoint.
func (s *Server) RegisterHealthHandler(f HealthFunc) {
	s.healthHandler.setHandler(f)
}

func (s *Server) onError(err error) {
	// TODO(prashant): Expose incoming call errors through options for NewServer.
	s.log.Errorf("thrift Server error: %v", err)
}

func (s *Server) handleStandard(origCtx context.Context, handler handler, method string, call *tchannel.InboundCall) error {
	reader, err := call.Arg2Reader()
	if err != nil {
		return err
	}
	headers, err := readHeaders(reader)
	if err != nil {
		return err
	}
	if err := reader.Close(); err != nil {
		return err
	}

	reader, err = call.Arg3Reader()
	if err != nil {
		return err
	}

	ctx := WithHeaders(origCtx, headers)
	protocol := thrift.NewTBinaryProtocolTransport(&readWriterTransport{Reader: reader})
	success, resp, err := handler.standard.Handle(ctx, method, protocol)
	if err != nil {
		reader.Close()
		call.Response().SendSystemError(err)
		return nil
	}
	if err := reader.Close(); err != nil {
		return err
	}

	if !success {
		call.Response().SetApplicationError()
	}

	writer, err := call.Response().Arg2Writer()
	if err != nil {
		return err
	}

	if err := writeHeaders(writer, ctx.ResponseHeaders()); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}

	writer, err = call.Response().Arg3Writer()
	protocol = thrift.NewTBinaryProtocolTransport(&readWriterTransport{Writer: writer})
	resp.Write(protocol)
	err = writer.Close()

	if handler.postResponseCB != nil {
		handler.postResponseCB(method, resp)
	}

	return err
}

// What if the arg3 is
// flags:4

// if flags && streamed
// then rest is chunk~4 chunk~4
// else it's just
// thrift struct.

// Make a new writer that will on first write do the writing of headers and then opens the arg3 writer
// if !written {
//   write arg2.
//   written = true
// }

// StreamWriter is the interface for writing multiple Thrift structs.
type StreamWriter interface {
	Flush() error
	Write(resp thrift.TStruct) error
	Close(err error) error
}

type streamWriter struct {
	call        *tchannel.InboundCall
	protocol    thrift.TProtocol
	arg2Written bool
}

func (w *streamWriter) sendArg2() error {

	w.arg2Written = true
	return nil
}

func (w *streamWriter) Flush() error {
	if !w.arg2Written {
		if err := w.sendArg2(); err != nil {
			return err
		}
	}

	return w.Flush()
}

func (w *streamWriter) Write(resp thrift.TStruct) error {
	if !w.arg2Written {
		if err := w.sendArg2(); err != nil {
			return err
		}
	}

	// Write out the struct
	if err := resp.Write(w.protocol); err != nil {
		return err
	}

	return nil
}

func (w *streamWriter) Close(err error) error {
	if err != nil && w.arg2Written {
		// All errors must be sent by now
		return fmt.Errorf("errors must be sent before any stream results")
	}

	if err != nil {
		w.call.Response().SetApplicationError()

		// NOT STREAMING?
	}

	if !w.arg2Written {
		if err := w.sendArg2(); err != nil {
			return err
		}

	}

	_, err = w.call.Response().Arg3Writer()
	if err != nil {
		return err
	}

	w.protocol.Flush()
	return nil
}

// Each language can then just translate directly?

// Handle handles an incoming TChannel call and forwards it to the correct handler.
func (s *Server) Handle(ctx context.Context, call *tchannel.InboundCall) {
	parts := strings.Split(string(call.Operation()), "::")
	if len(parts) != 2 {
		log.Fatalf("Handle got call for %v which does not match the expected call format", parts)
	}

	service, method := parts[0], parts[1]
	s.mut.RLock()
	handler, ok := s.handlers[service]
	s.mut.RUnlock()
	if !ok {
		log.Fatalf("Handle got call for service %v which is not registered", service)
	}

	// TODO(prashant): Logic for reading headers should not be duplicated.
	if handler.standard != nil {
		if err := s.handleStandard(ctx, handler, method, call); err != nil {
			s.onError(err)
		}
	} else {
		reader, err := call.Arg2Reader()
		if err != nil {
			log.Fatal(err)
		}
		headers, err := readHeaders(reader)
		if err != nil {
			log.Fatal(err)
		}
		if err := reader.Close(); err != nil {
			log.Fatal(err)
		}
		ctx := WithHeaders(ctx, headers)
		// TODO(prashant): Call post-response callback for each response write?
		handler.streaming.HandleStreaming(ctx, call)
	}
}

func (s *Server) HandleStreaming(ctx context.Context, call *tchannel.InboundCall) {
	parts := strings.Split(string(call.Operation()), "::")
	if len(parts) != 2 {
		log.Fatalf("Handle got call for %v which does not match the expected call format", parts)
	}

	service := parts[0]
	s.mut.RLock()
	handler, ok := s.handlers[service]
	s.mut.RUnlock()
	if !ok {
		log.Fatalf("Handle got call for service %v which is not registered", service)
	}

	reader, err := call.Arg2Reader()
	if err != nil {
		log.Fatal(err)
	}
	headers, err := readHeaders(reader)
	if err != nil {
		log.Fatal(err)
	}
	if err := reader.Close(); err != nil {
		log.Fatal(err)
	}
	tctx := WithHeaders(ctx, headers)
	// TODO(prashant): Call post-response callback for each response write?
	handler.streaming.HandleStreaming(tctx, call)
}
