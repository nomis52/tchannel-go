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

package thrift_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	// Test is in a separate package to avoid circular dependencies.
	. "github.com/uber/tchannel-go/thrift"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	tchannel "github.com/uber/tchannel-go"
	"github.com/uber/tchannel-go/testutils"
	gen "github.com/uber/tchannel-go/thrift/gen-go/test"
	"github.com/uber/tchannel-go/thrift/mocks"
)

// Generate the service mocks using go generate.
//go:generate mockery -name TChanSimpleService
//go:generate mockery -name TChanSecondService

type testArgs struct {
	server *Server
	s1     *mocks.TChanSimpleService
	s2     *mocks.TChanSecondService
	c1     gen.TChanSimpleService
	c2     gen.TChanSecondService
}

func ctxArg() mock.AnythingOfTypeArgument {
	return mock.AnythingOfType("*tchannel.headerCtx")
}

func TestThriftArgs(t *testing.T) {
	withSetup(t, func(ctx Context, args testArgs) {
		arg := &gen.Data{
			B1: true,
			S2: "str",
			I3: 102,
		}
		ret := &gen.Data{
			B1: false,
			S2: "return-str",
			I3: 105,
		}

		args.s1.On("Call", ctxArg(), arg).Return(ret, nil)
		got, err := args.c1.Call(ctx, arg)
		require.NoError(t, err)
		assert.Equal(t, ret, got)
	})
}

func TestRequest(t *testing.T) {
	withSetup(t, func(ctx Context, args testArgs) {
		args.s1.On("Simple", ctxArg()).Return(nil)
		require.NoError(t, args.c1.Simple(ctx))
	})
}

func TestRetryRequest(t *testing.T) {
	withSetup(t, func(ctx Context, args testArgs) {
		count := 0
		args.s1.On("Simple", ctxArg()).Return(tchannel.ErrServerBusy).
			Run(func(args mock.Arguments) {
			count++
		})
		require.Error(t, args.c1.Simple(ctx), "Simple expected to fail")
		assert.Equal(t, 5, count, "Expected Simple to be retried 5 times")
	})
}

func TestRequestSubChannel(t *testing.T) {
	ctx, cancel := NewContext(time.Second)
	defer cancel()

	tchan, err := tchannel.NewChannel("svc1", nil)
	require.NoError(t, err, "server NewChannel failed")
	require.NoError(t, tchan.ListenAndServe(":0"), "Listen failed")
	defer tchan.Close()

	clientCh, err := tchannel.NewChannel("client", nil)
	require.NoError(t, err, "client NewChannel failed")
	defer clientCh.Close()
	clientCh.Peers().Add(tchan.PeerInfo().HostPort)

	tests := []tchannel.Registrar{tchan, tchan.GetSubChannel("svc2"), tchan.GetSubChannel("svc3")}
	for _, ch := range tests {
		mockHandler := new(mocks.TChanSecondService)
		server := NewServer(ch)
		server.Register(gen.NewTChanSecondServiceServer(mockHandler))

		client := NewClient(clientCh, ch.ServiceName(), nil)
		secondClient := gen.NewTChanSecondServiceClient(client)

		echoArg := ch.ServiceName()
		echoRes := echoArg + "-echo"
		mockHandler.On("Echo", ctxArg(), echoArg).Return(echoRes, nil)
		res, err := secondClient.Echo(ctx, echoArg)
		assert.NoError(t, err, "Echo failed")
		assert.Equal(t, echoRes, res)
	}
}

func TestLargeRequest(t *testing.T) {
	arg := testutils.RandString(100000)
	res := strings.ToLower(arg)

	fmt.Println(len(arg))
	withSetup(t, func(ctx Context, args testArgs) {
		args.s2.On("Echo", ctxArg(), arg).Return(res, nil)

		got, err := args.c2.Echo(ctx, arg)
		if assert.NoError(t, err, "Echo got error") {
			assert.Equal(t, res, got, "Echo got unexpected response")
		}
	})
}

func TestThriftError(t *testing.T) {
	thriftErr := &gen.SimpleErr{
		Message: "this is the error",
	}
	withSetup(t, func(ctx Context, args testArgs) {
		args.s1.On("Simple", ctxArg()).Return(thriftErr)
		got := args.c1.Simple(ctx)
		require.Error(t, got)
		require.Equal(t, thriftErr, got)
	})
}

func TestUnknownError(t *testing.T) {
	withSetup(t, func(ctx Context, args testArgs) {
		args.s1.On("Simple", ctxArg()).Return(errors.New("unexpected err"))
		got := args.c1.Simple(ctx)
		require.Error(t, got)
		require.Equal(t, tchannel.NewSystemError(tchannel.ErrCodeUnexpected, "unexpected err"), got)
	})
}

func TestMultiple(t *testing.T) {
	withSetup(t, func(ctx Context, args testArgs) {
		args.s1.On("Simple", ctxArg()).Return(nil)
		args.s2.On("Echo", ctxArg(), "test1").Return("test2", nil)

		require.NoError(t, args.c1.Simple(ctx))
		res, err := args.c2.Echo(ctx, "test1")
		require.NoError(t, err)
		require.Equal(t, "test2", res)
	})
}

func TestHeaders(t *testing.T) {
	reqHeaders := map[string]string{"header1": "value1", "header2": "value2"}
	respHeaders := map[string]string{"resp1": "value1-resp", "resp2": "value2-resp"}

	withSetup(t, func(ctx Context, args testArgs) {
		args.s1.On("Simple", ctxArg()).Return(nil).Run(func(args mock.Arguments) {
			ctx := args.Get(0).(Context)
			assert.Equal(t, reqHeaders, ctx.Headers(), "request headers mismatch")
			ctx.SetResponseHeaders(respHeaders)
		})

		ctx = WithHeaders(ctx, reqHeaders)
		require.NoError(t, args.c1.Simple(ctx))
		assert.Equal(t, respHeaders, ctx.ResponseHeaders(), "response headers mismatch")
	})
}

func TestClientHostPort(t *testing.T) {
	ctx, cancel := NewContext(time.Second)
	defer cancel()

	s1ch := testutils.NewServer(t, nil)
	s2ch := testutils.NewServer(t, nil)
	defer s1ch.Close()
	defer s2ch.Close()

	s1ch.Peers().Add(s2ch.PeerInfo().HostPort)
	s2ch.Peers().Add(s1ch.PeerInfo().HostPort)

	mock1, mock2 := new(mocks.TChanSecondService), new(mocks.TChanSecondService)
	NewServer(s1ch).Register(gen.NewTChanSecondServiceServer(mock1))
	NewServer(s2ch).Register(gen.NewTChanSecondServiceServer(mock2))

	// When we call using a normal client, it can only call the other server (only peer).
	c1 := gen.NewTChanSecondServiceClient(NewClient(s1ch, s2ch.PeerInfo().ServiceName, nil))
	mock2.On("Echo", ctxArg(), "call1").Return("call1", nil)
	res, err := c1.Echo(ctx, "call1")
	assert.NoError(t, err, "call1 failed")
	assert.Equal(t, "call1", res)

	// When we call using a client that specifies host:port, it should call that server.
	c2 := gen.NewTChanSecondServiceClient(NewClient(s1ch, s1ch.PeerInfo().ServiceName, &ClientOptions{
		HostPort: s1ch.PeerInfo().HostPort,
	}))
	mock1.On("Echo", ctxArg(), "call2").Return("call2", nil)
	res, err = c2.Echo(ctx, "call2")
	assert.NoError(t, err, "call2 failed")
	assert.Equal(t, "call2", res)
}

func TestRegisterPostResponseCB(t *testing.T) {
	withSetup(t, func(ctx Context, args testArgs) {
		arg := &gen.Data{
			B1: true,
			S2: "str",
			I3: 102,
		}
		ret := &gen.Data{
			B1: false,
			S2: "return-str",
			I3: 105,
		}

		called := make(chan struct{})
		cb := func(method string, response thrift.TStruct) {
			assert.Equal(t, "Call", method)
			res, ok := response.(*gen.SimpleServiceCallResult)
			if assert.True(t, ok, "response type should be Result struct") {
				assert.Equal(t, ret, res.GetSuccess(), "result should be returned value")
			}
			close(called)
		}
		args.server.Register(gen.NewTChanSimpleServiceServer(args.s1), OptPostResponse(cb))

		args.s1.On("Call", ctxArg(), arg).Return(ret, nil)
		res, err := args.c1.Call(ctx, arg)
		require.NoError(t, err, "Call failed")
		assert.Equal(t, res, ret, "Call return value wrong")
		select {
		case <-time.After(time.Second):
			t.Errorf("post-response callback not called")
		case <-called:
		}
	})
}

func TestThriftTimeout(t *testing.T) {
	defer testutils.SetTimeout(t, time.Second)()

	withSetup(t, func(ctx Context, args testArgs) {
		handler := make(chan struct{})

		args.s2.On("Echo", ctxArg(), "asd").Return("asd", nil).Run(func(args mock.Arguments) {
			time.Sleep(15 * time.Millisecond)
			close(handler)
			return
		})

		ctx, cancel := NewContext(10 * time.Millisecond)
		defer cancel()

		_, err := args.c2.Echo(ctx, "asd")
		assert.Equal(t, err, tchannel.ErrTimeout, "Expect call to time out")

		// Wait for the handler to return, otherwise the test ends before the Server gets an error.
		<-handler
	})
}

func withSetup(t *testing.T, f func(ctx Context, args testArgs)) {
	args := testArgs{
		s1: new(mocks.TChanSimpleService),
		s2: new(mocks.TChanSecondService),
	}

	ctx, cancel := NewContext(time.Second)
	defer cancel()

	// Start server
	ch, server := setupServer(t, args.s1, args.s2)
	defer ch.Close()
	args.server = server

	// Get client1
	args.c1, args.c2 = getClients(t, ch)

	f(ctx, args)

	args.s1.AssertExpectations(t)
	args.s2.AssertExpectations(t)
}

func setupServer(t *testing.T, h *mocks.TChanSimpleService, sh *mocks.TChanSecondService) (*tchannel.Channel, *Server) {
	ch := testutils.NewServer(t, nil)
	server := NewServer(ch)
	server.Register(gen.NewTChanSimpleServiceServer(h))
	server.Register(gen.NewTChanSecondServiceServer(sh))
	return ch, server
}

func getClients(t *testing.T, serverCh *tchannel.Channel) (gen.TChanSimpleService, gen.TChanSecondService) {
	serverInfo := serverCh.PeerInfo()
	ch := testutils.NewClient(t, nil)

	ch.Peers().Add(serverInfo.HostPort)
	client := NewClient(ch, serverInfo.ServiceName, nil)

	simpleClient := gen.NewTChanSimpleServiceClient(client)
	secondClient := gen.NewTChanSecondServiceClient(client)
	return simpleClient, secondClient
}
