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
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/uber/tchannel-go"
)

// client implements TChanStreamingClient and makes outgoing Thrift calls.
type client struct {
	tcommon

	sc          *tchannel.SubChannel
	serviceName string
	opts        ClientOptions
}

// ClientOptions are options to customize the client.
type ClientOptions struct {
	// HostPort specifies a specific server to hit.
	HostPort string
}

// NewClient returns a Client that makes calls over the given tchannel to the given Hyperbahn service.
func NewClient(ch *tchannel.Channel, serviceName string, opts *ClientOptions) TChanStreamingClient {
	client := &client{
		sc:          ch.GetSubChannel(serviceName),
		serviceName: serviceName,
	}
	if opts != nil {
		client.opts = *opts
	}
	return client
}

func (c *client) StartCall(ctx Context, fullMethod string) (*tchannel.OutboundCall, tchannel.ArgWriter, error) {
	call, err := c.sc.BeginCall(ctx, fullMethod, &tchannel.CallOptions{Format: tchannel.Thrift})
	if err != nil {
		return nil, nil, err
	}

	writer, err := call.Arg2Writer()
	if err != nil {
		return nil, nil, err
	}
	if err := writeHeaders(writer, ctx.Headers()); err != nil {
		return nil, nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, nil, err
	}

	writer, err = call.Arg3Writer()
	if err != nil {
		return nil, nil, err
	}

	return call, writer, nil
}

func (c *client) Call(ctx Context, thriftService, methodName string, req, resp thrift.TStruct) (bool, error) {
	var peer *tchannel.Peer
	if c.opts.HostPort != "" {
		peer = c.sc.Peers().GetOrAdd(c.opts.HostPort)
	} else {
		var err error
		peer, err = c.sc.Peers().Get()
		if err != nil {
			return false, err
		}
	}
	call, err := peer.BeginCall(ctx, c.serviceName, thriftService+"::"+methodName, &tchannel.CallOptions{Format: tchannel.Thrift})
	if err != nil {
		return false, err
	}

	writer, err := call.Arg2Writer()
	if err != nil {
		return false, err
	}
	if err := writeHeaders(writer, ctx.Headers()); err != nil {
		return false, err
	}
	if err := writer.Close(); err != nil {
		return false, err
	}

	writer, err = call.Arg3Writer()
	if err != nil {
		return false, err
	}

	protocol := thrift.NewTBinaryProtocolTransport(&readWriterTransport{Writer: writer})
	if err := req.Write(protocol); err != nil {
		return false, err
	}
	if err := writer.Close(); err != nil {
		return false, err
	}

	reader, err := call.Response().Arg2Reader()
	if err != nil {
		return false, err
	}

	headers, err := readHeaders(reader)
	if err != nil {
		return false, err
	}
	ctx.SetResponseHeaders(headers)
	if err := reader.Close(); err != nil {
		return false, err
	}

	success := !call.Response().ApplicationError()
	reader, err = call.Response().Arg3Reader()
	if err != nil {
		return success, err
	}

	protocol = thrift.NewTBinaryProtocolTransport(&readWriterTransport{Reader: reader})
	if err := resp.Read(protocol); err != nil {
		return success, err
	}
	if err := reader.Close(); err != nil {
		return success, err
	}

	return success, nil
}
