// Autogenerated by Thrift Compiler (1.0.0-dev)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package stream

import (
	"bytes"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
)

// (needed to ensure safety because of naive import list construction.)
var _ = thrift.ZERO
var _ = fmt.Printf
var _ = bytes.Equal

type TestStream2 interface {
	TestStream

	// Parameters:
	//  - Prefix
	OutStream2(prefix string) (r *SStringStream, err error)
}

type TestStream2Client struct {
	*TestStreamClient
}

func NewTestStream2ClientFactory(t thrift.TTransport, f thrift.TProtocolFactory) *TestStream2Client {
	return &TestStream2Client{TestStreamClient: NewTestStreamClientFactory(t, f)}
}

func NewTestStream2ClientProtocol(t thrift.TTransport, iprot thrift.TProtocol, oprot thrift.TProtocol) *TestStream2Client {
	return &TestStream2Client{TestStreamClient: NewTestStreamClientProtocol(t, iprot, oprot)}
}

// Parameters:
//  - Prefix
func (p *TestStream2Client) OutStream2(prefix string) (r *SStringStream, err error) {
	if err = p.sendOutStream2(prefix); err != nil {
		return
	}
	return p.recvOutStream2()
}

func (p *TestStream2Client) sendOutStream2(prefix string) (err error) {
	oprot := p.OutputProtocol
	if oprot == nil {
		oprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.OutputProtocol = oprot
	}
	p.SeqId++
	if err = oprot.WriteMessageBegin("OutStream2", thrift.CALL, p.SeqId); err != nil {
		return
	}
	args := TestStream2OutStream2Args{
		Prefix: prefix,
	}
	if err = args.Write(oprot); err != nil {
		return
	}
	if err = oprot.WriteMessageEnd(); err != nil {
		return
	}
	return oprot.Flush()
}

func (p *TestStream2Client) recvOutStream2() (value *SStringStream, err error) {
	iprot := p.InputProtocol
	if iprot == nil {
		iprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.InputProtocol = iprot
	}
	method, mTypeId, seqId, err := iprot.ReadMessageBegin()
	if err != nil {
		return
	}
	if method != "OutStream2" {
		err = thrift.NewTApplicationException(thrift.WRONG_METHOD_NAME, "OutStream2 failed: wrong method name")
		return
	}
	if p.SeqId != seqId {
		err = thrift.NewTApplicationException(thrift.BAD_SEQUENCE_ID, "OutStream2 failed: out of sequence response")
		return
	}
	if mTypeId == thrift.EXCEPTION {
		error13 := thrift.NewTApplicationException(thrift.UNKNOWN_APPLICATION_EXCEPTION, "Unknown Exception")
		var error14 error
		error14, err = error13.Read(iprot)
		if err != nil {
			return
		}
		if err = iprot.ReadMessageEnd(); err != nil {
			return
		}
		err = error14
		return
	}
	if mTypeId != thrift.REPLY {
		err = thrift.NewTApplicationException(thrift.INVALID_MESSAGE_TYPE_EXCEPTION, "OutStream2 failed: invalid message type")
		return
	}
	result := TestStream2OutStream2Result{}
	if err = result.Read(iprot); err != nil {
		return
	}
	if err = iprot.ReadMessageEnd(); err != nil {
		return
	}
	value = result.GetSuccess()
	return
}

type TestStream2Processor struct {
	*TestStreamProcessor
}

func NewTestStream2Processor(handler TestStream2) *TestStream2Processor {
	self15 := &TestStream2Processor{NewTestStreamProcessor(handler)}
	self15.AddToProcessorMap("OutStream2", &testStream2ProcessorOutStream2{handler: handler})
	return self15
}

type testStream2ProcessorOutStream2 struct {
	handler TestStream2
}

func (p *testStream2ProcessorOutStream2) Process(seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	args := TestStream2OutStream2Args{}
	if err = args.Read(iprot); err != nil {
		iprot.ReadMessageEnd()
		x := thrift.NewTApplicationException(thrift.PROTOCOL_ERROR, err.Error())
		oprot.WriteMessageBegin("OutStream2", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush()
		return false, err
	}

	iprot.ReadMessageEnd()
	result := TestStream2OutStream2Result{}
	var retval *SStringStream
	var err2 error
	if retval, err2 = p.handler.OutStream2(args.Prefix); err2 != nil {
		x := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, "Internal error processing OutStream2: "+err2.Error())
		oprot.WriteMessageBegin("OutStream2", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush()
		return true, err2
	} else {
		result.Success = retval
	}
	if err2 = oprot.WriteMessageBegin("OutStream2", thrift.REPLY, seqId); err2 != nil {
		err = err2
	}
	if err2 = result.Write(oprot); err == nil && err2 != nil {
		err = err2
	}
	if err2 = oprot.WriteMessageEnd(); err == nil && err2 != nil {
		err = err2
	}
	if err2 = oprot.Flush(); err == nil && err2 != nil {
		err = err2
	}
	if err != nil {
		return
	}
	return true, err
}

// HELPER FUNCTIONS AND STRUCTURES

// Attributes:
//  - Prefix
type TestStream2OutStream2Args struct {
	Prefix string `thrift:"prefix,1" json:"prefix"`
}

func NewTestStream2OutStream2Args() *TestStream2OutStream2Args {
	return &TestStream2OutStream2Args{}
}

func (p *TestStream2OutStream2Args) GetPrefix() string {
	return p.Prefix
}
func (p *TestStream2OutStream2Args) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.readField1(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *TestStream2OutStream2Args) readField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		p.Prefix = v
	}
	return nil
}

func (p *TestStream2OutStream2Args) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("OutStream2_args"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *TestStream2OutStream2Args) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("prefix", thrift.STRING, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:prefix: ", p), err)
	}
	if err := oprot.WriteString(string(p.Prefix)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.prefix (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:prefix: ", p), err)
	}
	return err
}

func (p *TestStream2OutStream2Args) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("TestStream2OutStream2Args(%+v)", *p)
}

// Attributes:
//  - Success
type TestStream2OutStream2Result struct {
	Success *SStringStream `thrift:"success,0" json:"success,omitempty"`
}

func NewTestStream2OutStream2Result() *TestStream2OutStream2Result {
	return &TestStream2OutStream2Result{}
}

var TestStream2OutStream2Result_Success_DEFAULT *SStringStream

func (p *TestStream2OutStream2Result) GetSuccess() *SStringStream {
	if !p.IsSetSuccess() {
		return TestStream2OutStream2Result_Success_DEFAULT
	}
	return p.Success
}
func (p *TestStream2OutStream2Result) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *TestStream2OutStream2Result) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 0:
			if err := p.readField0(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *TestStream2OutStream2Result) readField0(iprot thrift.TProtocol) error {
	p.Success = &SStringStream{}
	if err := p.Success.Read(iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.Success), err)
	}
	return nil
}

func (p *TestStream2OutStream2Result) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("OutStream2_result"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if err := p.writeField0(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *TestStream2OutStream2Result) writeField0(oprot thrift.TProtocol) (err error) {
	if p.IsSetSuccess() {
		if err := oprot.WriteFieldBegin("success", thrift.STRUCT, 0); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 0:success: ", p), err)
		}
		if err := p.Success.Write(oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.Success), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 0:success: ", p), err)
		}
	}
	return err
}

func (p *TestStream2OutStream2Result) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("TestStream2OutStream2Result(%+v)", *p)
}
