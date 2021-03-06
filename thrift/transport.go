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
	"io"
	"sync"

	"github.com/apache/thrift/lib/go/thrift"
)

// readerWriterTransport is a transport that reads and writes from the underlying Reader/Writer.
type readWriterTransport struct {
	io.Writer
	io.Reader
}

func (t *readWriterTransport) Open() error {
	return nil
}

func (t *readWriterTransport) Flush() error {
	return nil
}

func (t *readWriterTransport) IsOpen() bool {
	return true
}

func (t *readWriterTransport) Close() error {
	return nil
}

// RemainingBytes returns the max number of bytes (same as Thrift's StreamTransport) as we
// do not know how many bytes we have left.
func (t *readWriterTransport) RemainingBytes() uint64 {
	const maxSize = ^uint64(0)
	return maxSize
}

var _ thrift.TTransport = &readWriterTransport{}

type thriftProtocol struct {
	transport *readWriterTransport
	protocol  *thrift.TBinaryProtocol
}

var thriftProtocolPool = sync.Pool{
	New: func() interface{} {
		transport := &readWriterTransport{}
		protocol := thrift.NewTBinaryProtocolTransport(transport)
		return &thriftProtocol{transport, protocol}
	},
}

func getProtocolWriter(writer io.Writer) *thriftProtocol {
	wp := thriftProtocolPool.Get().(*thriftProtocol)
	wp.transport.Reader = nil
	wp.transport.Writer = writer
	return wp
}

func getProtocolReader(reader io.Reader) *thriftProtocol {
	wp := thriftProtocolPool.Get().(*thriftProtocol)
	wp.transport.Reader = reader
	wp.transport.Writer = nil
	return wp
}
