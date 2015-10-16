// Autogenerated. Code generated by thrift-gen. Do not modify.
package stream

import (
	"fmt"
	"io"

	athrift "github.com/apache/thrift/lib/go/thrift"
	"github.com/uber/tchannel-go"
	"github.com/uber/tchannel-go/thrift"
)

// Used to avoid unused warnings for non-streaming services.
var _ = tchannel.NewChannel
var _ = io.Reader(nil)

// Interfaces for the service and client for the services defined in the IDL.

type TChanUniqC interface {
}

// TChanUniqCServer is the interface that must be implemented by a handler.
type TChanUniqCServer interface {
	Run(ctx thrift.Context, call *UniqCRunInCall) error
}

// TChanUniqCClient is the interface is used to make remote calls.
type TChanUniqCClient interface {
	Run(ctx thrift.Context) (*UniqCRunOutCall, error)
}

type TChanUniqC2 interface {
	TChanUniqC
}

// TChanUniqC2Server is the interface that must be implemented by a handler.
type TChanUniqC2Server interface {
	TChanUniqCServer

	Fakerun(ctx thrift.Context, call *UniqC2FakerunInCall) error
}

// TChanUniqC2Client is the interface is used to make remote calls.
type TChanUniqC2Client interface {
	TChanUniqCClient

	Fakerun(ctx thrift.Context) (*UniqC2FakerunOutCall, error)
}

// Implementation of a client and service handler.

type tchanUniqCClient struct {
	thriftService string
	client        thrift.TChanStreamingClient
}

func newTChanUniqCClient(thriftService string, client thrift.TChanStreamingClient) *tchanUniqCClient {
	return &tchanUniqCClient{
		thriftService,
		client,
	}
}

func NewTChanUniqCClient(client thrift.TChanStreamingClient) TChanUniqCClient {
	return newTChanUniqCClient("UniqC", client)
}

func (c *tchanUniqCClient) Run(ctx thrift.Context) (*UniqCRunOutCall, error) {
	call, writer, err := c.client.StartCall(ctx, "UniqC::run")
	if err != nil {
		return nil, err
	}

	outCall := &UniqCRunOutCall{
		c:    c.client,
		call: call,
	}

	outCall.writer = writer

	return outCall, nil
}

type tchanUniqCServer struct {
	handler TChanUniqCServer
	common  thrift.TCommon
}

func newTChanUniqCServer(handler TChanUniqCServer) *tchanUniqCServer {
	return &tchanUniqCServer{
		handler,
		nil, /* common */
	}
}

func NewTChanUniqCServer(handler TChanUniqCServer) thrift.TChanStreamingServer {
	return newTChanUniqCServer(handler)
}

func (s *tchanUniqCServer) Service() string {
	return "UniqC"
}

func (s *tchanUniqCServer) SetCommon(common thrift.TCommon) {
	s.common = common
}

func (s *tchanUniqCServer) Methods() []string {
	return []string{}
}

func (s *tchanUniqCServer) StreamingMethods() []string {
	return []string{
		"run",
	}
}

func (s *tchanUniqCServer) HandleStreaming(ctx thrift.Context, call *tchannel.InboundCall) error {
	methodName := string(call.Operation())
	arg3Reader, err := call.Arg3Reader()
	if err != nil {
		return err
	}
	switch methodName {
	case "UniqC::run":
		return s.handleRun(ctx, call, arg3Reader)
	}
	return fmt.Errorf("method %v not found in service %v", methodName, s.Service())
}

func (s *tchanUniqCServer) Handle(ctx thrift.Context, methodName string, protocol athrift.TProtocol) (bool, athrift.TStruct, error) {
	switch methodName {

	default:
		return false, nil, fmt.Errorf("method %v not found in service %v", methodName, s.Service())
	}
}

func (s *tchanUniqCServer) handleRun(ctx thrift.Context, tcall *tchannel.InboundCall, arg3Reader io.ReadCloser) error {
	call := &UniqCRunInCall{
		c:    s.common,
		call: tcall,
		ctx:  ctx,
	}

	call.reader = arg3Reader

	err :=
		s.handler.Run(ctx, call)
	if err != nil {
		// TODO: encode any Thrift exceptions here.
		return err
	}

	if err := call.checkWriter(); err != nil {
		return err
	}

	// TODO: we may want to Close the writer if it's not already closed.

	return nil
}

// UniqCRunInCall is the object used to stream arguments and write
// response headers for incoming calls.
type UniqCRunInCall struct {
	c    thrift.TCommon
	call *tchannel.InboundCall
	ctx  thrift.Context

	reader io.ReadCloser

	writer tchannel.ArgWriter
}

// Read returns the next argument, if any is available. If there are no more
// arguments left, it will return io.EOF.
func (c *UniqCRunInCall) Read() (*String, error) {
	var req String
	if err := c.c.ReadStreamStruct(c.reader, func(protocol athrift.TProtocol) error {
		return req.Read(protocol)
	}); err != nil {
		return nil, err
	}

	return &req, nil
}

// SetResponseHeaders sets the response headers. This must be called before any
// streaming responses are sent.
func (c *UniqCRunInCall) SetResponseHeaders(headers map[string]string) error {
	if c.writer != nil {
		// arg3 is already being written, headers must be set first
		return fmt.Errorf("cannot set headers after writing streaming responses")
	}

	c.ctx.SetResponseHeaders(headers)
	return nil
}

func (c *UniqCRunInCall) writeResponseHeaders() error {
	if c.writer != nil {
		// arg3 is already being written, headers must be set first
		return fmt.Errorf("cannot set headers after writing streaming responses")
	}

	// arg2 writer should be used to write headers
	arg2Writer, err := c.call.Response().Arg2Writer()
	if err != nil {
		return err
	}

	headers := c.ctx.ResponseHeaders()
	if err := c.c.WriteHeaders(arg2Writer, headers); err != nil {
		return err
	}

	return arg2Writer.Close()
}

// checkWriter creates the arg3 writer if it has not been created.
// Before the arg3 writer is created, response headers are sent.
func (c *UniqCRunInCall) checkWriter() error {
	if c.writer == nil {
		if err := c.writeResponseHeaders(); err != nil {
			return err
		}

		writer, err := c.call.Response().Arg3Writer()
		if err != nil {
			return err
		}
		c.writer = writer
	}
	return nil
}

// Write writes a result to the response stream. The written items may not
// be sent till Flush or Done is called.
func (c *UniqCRunInCall) Write(arg *SCount) error {
	if err := c.checkWriter(); err != nil {
		return err
	}
	return c.c.WriteStreamStruct(c.writer, arg)
}

// Flush flushes headers (if they have not yet been sent) and any written results.
func (c *UniqCRunInCall) Flush() error {
	if err := c.checkWriter(); err != nil {
		return err
	}
	return c.writer.Flush()
}

// Done closes the response stream and should be called after all results have been written.
func (c *UniqCRunInCall) Done() error {
	if err := c.checkWriter(); err != nil {
		return err
	}
	return c.writer.Close()
}

// UniqCRunOutCall is the object used to stream arguments/results and
// read response headers for outgoing calls.
type UniqCRunOutCall struct {
	c               thrift.TCommon
	call            *tchannel.OutboundCall
	responseHeaders map[string]string
	reader          io.ReadCloser
	writer          tchannel.ArgWriter
}

// Write writes an argument to the request stream. The written items may not
// be sent till Flush or Done is called.
func (c *UniqCRunOutCall) Write(arg *String) error {
	return c.c.WriteStreamStruct(c.writer, arg)
}

// Flush flushes all written arguments.
func (c *UniqCRunOutCall) Flush() error {
	return c.writer.Flush()
}

// Done closes the request stream and should be called after all arguments have been written.
func (c *UniqCRunOutCall) Done() error {
	if err := c.writer.Close(); err != nil {
		return err
	}

	return nil
}

func (c *UniqCRunOutCall) checkReader() error {
	if c.reader == nil {
		arg2Reader, err := c.call.Response().Arg2Reader()
		if err != nil {
			return err
		}

		c.responseHeaders, err = c.c.ReadHeaders(arg2Reader)
		if err != nil {
			return err
		}
		if err := arg2Reader.Close(); err != nil {
			return err
		}

		reader, err := c.call.Response().Arg3Reader()
		if err != nil {
			return err
		}

		c.reader = reader
	}
	return nil
}

// Read returns the next result, if any is available. If there are no more
// results left, it will return io.EOF.
func (c *UniqCRunOutCall) Read() (*SCount, error) {
	if err := c.checkReader(); err != nil {
		return nil, err
	}
	var res SCount
	if err := c.c.ReadStreamStruct(c.reader, func(protocol athrift.TProtocol) error {
		return res.Read(protocol)
	}); err != nil {
		return nil, err
	}

	return &res, nil
}

// ResponseHeaders returns the response headers sent from the server. This will
// block until server headers have been received.
func (c *UniqCRunOutCall) ResponseHeaders() (map[string]string, error) {
	if err := c.checkReader(); err != nil {
		return nil, err
	}
	return c.responseHeaders, nil
}

type tchanUniqC2Client struct {
	tchanUniqCClient

	thriftService string
	client        thrift.TChanStreamingClient
}

func newTChanUniqC2Client(thriftService string, client thrift.TChanStreamingClient) *tchanUniqC2Client {
	return &tchanUniqC2Client{
		*newTChanUniqCClient(thriftService, client),
		thriftService,
		client,
	}
}

func NewTChanUniqC2Client(client thrift.TChanStreamingClient) TChanUniqC2Client {
	return newTChanUniqC2Client("UniqC2", client)
}

func (c *tchanUniqC2Client) Fakerun(ctx thrift.Context) (*UniqC2FakerunOutCall, error) {
	call, writer, err := c.client.StartCall(ctx, "UniqC2::fakerun")
	if err != nil {
		return nil, err
	}

	outCall := &UniqC2FakerunOutCall{
		c:    c.client,
		call: call,
	}

	outCall.writer = writer

	return outCall, nil
}

type tchanUniqC2Server struct {
	tchanUniqCServer

	handler TChanUniqC2Server
	common  thrift.TCommon
}

func newTChanUniqC2Server(handler TChanUniqC2Server) *tchanUniqC2Server {
	return &tchanUniqC2Server{
		*newTChanUniqCServer(handler),
		handler,
		nil, /* common */
	}
}

func NewTChanUniqC2Server(handler TChanUniqC2Server) thrift.TChanStreamingServer {
	return newTChanUniqC2Server(handler)
}

func (s *tchanUniqC2Server) Service() string {
	return "UniqC2"
}

func (s *tchanUniqC2Server) SetCommon(common thrift.TCommon) {
	s.common = common
	s.tchanUniqCServer.SetCommon(common)
}

func (s *tchanUniqC2Server) Methods() []string {
	return []string{}
}

func (s *tchanUniqC2Server) StreamingMethods() []string {
	return []string{
		"fakerun",

		"run",
	}
}

func (s *tchanUniqC2Server) HandleStreaming(ctx thrift.Context, call *tchannel.InboundCall) error {
	methodName := string(call.Operation())
	arg3Reader, err := call.Arg3Reader()
	if err != nil {
		return err
	}
	switch methodName {
	case "UniqC2::fakerun":
		return s.handleFakerun(ctx, call, arg3Reader)
	}
	return fmt.Errorf("method %v not found in service %v", methodName, s.Service())
}

func (s *tchanUniqC2Server) Handle(ctx thrift.Context, methodName string, protocol athrift.TProtocol) (bool, athrift.TStruct, error) {
	switch methodName {

	default:
		return false, nil, fmt.Errorf("method %v not found in service %v", methodName, s.Service())
	}
}

func (s *tchanUniqC2Server) handleFakerun(ctx thrift.Context, tcall *tchannel.InboundCall, arg3Reader io.ReadCloser) error {
	call := &UniqC2FakerunInCall{
		c:    s.common,
		call: tcall,
		ctx:  ctx,
	}

	call.reader = arg3Reader

	err :=
		s.handler.Fakerun(ctx, call)
	if err != nil {
		// TODO: encode any Thrift exceptions here.
		return err
	}

	if err := call.checkWriter(); err != nil {
		return err
	}

	// TODO: we may want to Close the writer if it's not already closed.

	return nil
}

// UniqC2FakerunInCall is the object used to stream arguments and write
// response headers for incoming calls.
type UniqC2FakerunInCall struct {
	c    thrift.TCommon
	call *tchannel.InboundCall
	ctx  thrift.Context

	reader io.ReadCloser

	writer tchannel.ArgWriter
}

// Read returns the next argument, if any is available. If there are no more
// arguments left, it will return io.EOF.
func (c *UniqC2FakerunInCall) Read() (*String, error) {
	var req String
	if err := c.c.ReadStreamStruct(c.reader, func(protocol athrift.TProtocol) error {
		return req.Read(protocol)
	}); err != nil {
		return nil, err
	}

	return &req, nil
}

// SetResponseHeaders sets the response headers. This must be called before any
// streaming responses are sent.
func (c *UniqC2FakerunInCall) SetResponseHeaders(headers map[string]string) error {
	if c.writer != nil {
		// arg3 is already being written, headers must be set first
		return fmt.Errorf("cannot set headers after writing streaming responses")
	}

	c.ctx.SetResponseHeaders(headers)
	return nil
}

func (c *UniqC2FakerunInCall) writeResponseHeaders() error {
	if c.writer != nil {
		// arg3 is already being written, headers must be set first
		return fmt.Errorf("cannot set headers after writing streaming responses")
	}

	// arg2 writer should be used to write headers
	arg2Writer, err := c.call.Response().Arg2Writer()
	if err != nil {
		return err
	}

	headers := c.ctx.ResponseHeaders()
	if err := c.c.WriteHeaders(arg2Writer, headers); err != nil {
		return err
	}

	return arg2Writer.Close()
}

// checkWriter creates the arg3 writer if it has not been created.
// Before the arg3 writer is created, response headers are sent.
func (c *UniqC2FakerunInCall) checkWriter() error {
	if c.writer == nil {
		if err := c.writeResponseHeaders(); err != nil {
			return err
		}

		writer, err := c.call.Response().Arg3Writer()
		if err != nil {
			return err
		}
		c.writer = writer
	}
	return nil
}

// Write writes a result to the response stream. The written items may not
// be sent till Flush or Done is called.
func (c *UniqC2FakerunInCall) Write(arg *SCount) error {
	if err := c.checkWriter(); err != nil {
		return err
	}
	return c.c.WriteStreamStruct(c.writer, arg)
}

// Flush flushes headers (if they have not yet been sent) and any written results.
func (c *UniqC2FakerunInCall) Flush() error {
	if err := c.checkWriter(); err != nil {
		return err
	}
	return c.writer.Flush()
}

// Done closes the response stream and should be called after all results have been written.
func (c *UniqC2FakerunInCall) Done() error {
	if err := c.checkWriter(); err != nil {
		return err
	}
	return c.writer.Close()
}

// UniqC2FakerunOutCall is the object used to stream arguments/results and
// read response headers for outgoing calls.
type UniqC2FakerunOutCall struct {
	c               thrift.TCommon
	call            *tchannel.OutboundCall
	responseHeaders map[string]string
	reader          io.ReadCloser
	writer          tchannel.ArgWriter
}

// Write writes an argument to the request stream. The written items may not
// be sent till Flush or Done is called.
func (c *UniqC2FakerunOutCall) Write(arg *String) error {
	return c.c.WriteStreamStruct(c.writer, arg)
}

// Flush flushes all written arguments.
func (c *UniqC2FakerunOutCall) Flush() error {
	return c.writer.Flush()
}

// Done closes the request stream and should be called after all arguments have been written.
func (c *UniqC2FakerunOutCall) Done() error {
	if err := c.writer.Close(); err != nil {
		return err
	}

	return nil
}

func (c *UniqC2FakerunOutCall) checkReader() error {
	if c.reader == nil {
		arg2Reader, err := c.call.Response().Arg2Reader()
		if err != nil {
			return err
		}

		c.responseHeaders, err = c.c.ReadHeaders(arg2Reader)
		if err != nil {
			return err
		}
		if err := arg2Reader.Close(); err != nil {
			return err
		}

		reader, err := c.call.Response().Arg3Reader()
		if err != nil {
			return err
		}

		c.reader = reader
	}
	return nil
}

// Read returns the next result, if any is available. If there are no more
// results left, it will return io.EOF.
func (c *UniqC2FakerunOutCall) Read() (*SCount, error) {
	if err := c.checkReader(); err != nil {
		return nil, err
	}
	var res SCount
	if err := c.c.ReadStreamStruct(c.reader, func(protocol athrift.TProtocol) error {
		return res.Read(protocol)
	}); err != nil {
		return nil, err
	}

	return &res, nil
}

// ResponseHeaders returns the response headers sent from the server. This will
// block until server headers have been received.
func (c *UniqC2FakerunOutCall) ResponseHeaders() (map[string]string, error) {
	if err := c.checkReader(); err != nil {
		return nil, err
	}
	return c.responseHeaders, nil
}
