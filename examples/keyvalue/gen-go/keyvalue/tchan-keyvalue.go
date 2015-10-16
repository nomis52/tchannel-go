// Autogenerated. Code generated by thrift-gen. Do not modify.
package keyvalue

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

type TChanAdmin interface {
	TChanBaseService

	ClearAll(ctx thrift.Context) error
}

// TChanAdminServer is the interface that must be implemented by a handler.
type TChanAdminServer interface {
	TChanBaseServiceServer

	ClearAll(ctx thrift.Context) error
}

// TChanAdminClient is the interface is used to make remote calls.
type TChanAdminClient interface {
	TChanBaseServiceClient

	ClearAll(ctx thrift.Context) error
}

type TChanKeyValue interface {
	TChanBaseService

	Get(ctx thrift.Context, key string) (string, error)
	Hi(ctx thrift.Context) error
	Set(ctx thrift.Context, key string, value string) error
}

// TChanKeyValueServer is the interface that must be implemented by a handler.
type TChanKeyValueServer interface {
	TChanBaseServiceServer

	Get(ctx thrift.Context, key string) (string, error)
	Hi(ctx thrift.Context) error
	Set(ctx thrift.Context, key string, value string) error
}

// TChanKeyValueClient is the interface is used to make remote calls.
type TChanKeyValueClient interface {
	TChanBaseServiceClient

	Get(ctx thrift.Context, key string) (string, error)
	Hi(ctx thrift.Context) error
	Set(ctx thrift.Context, key string, value string) error
}

type TChanBaseService interface {
	HealthCheck(ctx thrift.Context) (string, error)
}

// TChanBaseServiceServer is the interface that must be implemented by a handler.
type TChanBaseServiceServer interface {
	HealthCheck(ctx thrift.Context) (string, error)
}

// TChanBaseServiceClient is the interface is used to make remote calls.
type TChanBaseServiceClient interface {
	HealthCheck(ctx thrift.Context) (string, error)
}

// Implementation of a client and service handler.

type tchanAdminClient struct {
	tchanBaseServiceClient

	thriftService string
	client        thrift.TChanStreamingClient
}

func newTChanAdminClient(thriftService string, client thrift.TChanStreamingClient) *tchanAdminClient {
	return &tchanAdminClient{
		*newTChanBaseServiceClient(thriftService, client),
		thriftService,
		client,
	}
}

func NewTChanAdminClient(client thrift.TChanStreamingClient) TChanAdminClient {
	return newTChanAdminClient("Admin", client)
}

func (c *tchanAdminClient) ClearAll(ctx thrift.Context) error {
	var resp AdminClearAllResult
	args := AdminClearAllArgs{}
	success, err := c.client.Call(ctx, c.thriftService, "clearAll", &args, &resp)
	if err == nil && !success {
		if e := resp.NotAuthorized; e != nil {
			err = e
		}
	}

	return err
}

type tchanAdminServer struct {
	tchanBaseServiceServer

	handler TChanAdminServer
	common  thrift.TCommon
}

func newTChanAdminServer(handler TChanAdminServer) *tchanAdminServer {
	return &tchanAdminServer{
		*newTChanBaseServiceServer(handler),
		handler,
		nil, /* common */
	}
}

func NewTChanAdminServer(handler TChanAdminServer) thrift.TChanStreamingServer {
	return newTChanAdminServer(handler)
}

func (s *tchanAdminServer) Service() string {
	return "Admin"
}

func (s *tchanAdminServer) SetCommon(common thrift.TCommon) {
	s.common = common
	s.tchanBaseServiceServer.SetCommon(common)
}

func (s *tchanAdminServer) Methods() []string {
	return []string{
		"clearAll",

		"HealthCheck",
	}
}

func (s *tchanAdminServer) StreamingMethods() []string {
	return []string{}
}

func (s *tchanAdminServer) HandleStreaming(ctx thrift.Context, call *tchannel.InboundCall) error {
	methodName := string(call.Operation())
	return fmt.Errorf("method %v not found in service %v", methodName, s.Service())
}

func (s *tchanAdminServer) Handle(ctx thrift.Context, methodName string, protocol athrift.TProtocol) (bool, athrift.TStruct, error) {
	switch methodName {
	case "clearAll":
		return s.handleClearAll(ctx, protocol)

	case "HealthCheck":
		return s.handleHealthCheck(ctx, protocol)

	default:
		return false, nil, fmt.Errorf("method %v not found in service %v", methodName, s.Service())
	}
}

func (s *tchanAdminServer) handleClearAll(ctx thrift.Context, protocol athrift.TProtocol) (bool, athrift.TStruct, error) {
	var req AdminClearAllArgs
	var res AdminClearAllResult

	if err := req.Read(protocol); err != nil {
		return false, nil, err
	}

	err :=
		s.handler.ClearAll(ctx)

	if err != nil {
		switch v := err.(type) {
		case *NotAuthorized:
			res.NotAuthorized = v
		default:
			return false, nil, err
		}
	} else {
	}

	return err == nil, &res, nil
}

type tchanKeyValueClient struct {
	tchanBaseServiceClient

	thriftService string
	client        thrift.TChanStreamingClient
}

func newTChanKeyValueClient(thriftService string, client thrift.TChanStreamingClient) *tchanKeyValueClient {
	return &tchanKeyValueClient{
		*newTChanBaseServiceClient(thriftService, client),
		thriftService,
		client,
	}
}

func NewTChanKeyValueClient(client thrift.TChanStreamingClient) TChanKeyValueClient {
	return newTChanKeyValueClient("KeyValue", client)
}

func (c *tchanKeyValueClient) Get(ctx thrift.Context, key string) (string, error) {
	var resp KeyValueGetResult
	args := KeyValueGetArgs{
		Key: key,
	}
	success, err := c.client.Call(ctx, c.thriftService, "Get", &args, &resp)
	if err == nil && !success {
		if e := resp.NotFound; e != nil {
			err = e
		}
		if e := resp.InvalidKey; e != nil {
			err = e
		}
	}

	return resp.GetSuccess(), err
}

func (c *tchanKeyValueClient) Hi(ctx thrift.Context) error {
	var resp KeyValueHiResult
	args := KeyValueHiArgs{}
	success, err := c.client.Call(ctx, c.thriftService, "Hi", &args, &resp)
	if err == nil && !success {
	}

	return err
}

func (c *tchanKeyValueClient) Set(ctx thrift.Context, key string, value string) error {
	var resp KeyValueSetResult
	args := KeyValueSetArgs{
		Key:   key,
		Value: value,
	}
	success, err := c.client.Call(ctx, c.thriftService, "Set", &args, &resp)
	if err == nil && !success {
		if e := resp.InvalidKey; e != nil {
			err = e
		}
	}

	return err
}

type tchanKeyValueServer struct {
	tchanBaseServiceServer

	handler TChanKeyValueServer
	common  thrift.TCommon
}

func newTChanKeyValueServer(handler TChanKeyValueServer) *tchanKeyValueServer {
	return &tchanKeyValueServer{
		*newTChanBaseServiceServer(handler),
		handler,
		nil, /* common */
	}
}

func NewTChanKeyValueServer(handler TChanKeyValueServer) thrift.TChanStreamingServer {
	return newTChanKeyValueServer(handler)
}

func (s *tchanKeyValueServer) Service() string {
	return "KeyValue"
}

func (s *tchanKeyValueServer) SetCommon(common thrift.TCommon) {
	s.common = common
	s.tchanBaseServiceServer.SetCommon(common)
}

func (s *tchanKeyValueServer) Methods() []string {
	return []string{
		"Get",
		"Hi",
		"Set",

		"HealthCheck",
	}
}

func (s *tchanKeyValueServer) StreamingMethods() []string {
	return []string{}
}

func (s *tchanKeyValueServer) HandleStreaming(ctx thrift.Context, call *tchannel.InboundCall) error {
	methodName := string(call.Operation())
	return fmt.Errorf("method %v not found in service %v", methodName, s.Service())
}

func (s *tchanKeyValueServer) Handle(ctx thrift.Context, methodName string, protocol athrift.TProtocol) (bool, athrift.TStruct, error) {
	switch methodName {
	case "Get":
		return s.handleGet(ctx, protocol)
	case "Hi":
		return s.handleHi(ctx, protocol)
	case "Set":
		return s.handleSet(ctx, protocol)

	case "HealthCheck":
		return s.handleHealthCheck(ctx, protocol)

	default:
		return false, nil, fmt.Errorf("method %v not found in service %v", methodName, s.Service())
	}
}

func (s *tchanKeyValueServer) handleGet(ctx thrift.Context, protocol athrift.TProtocol) (bool, athrift.TStruct, error) {
	var req KeyValueGetArgs
	var res KeyValueGetResult

	if err := req.Read(protocol); err != nil {
		return false, nil, err
	}

	r, err :=
		s.handler.Get(ctx, req.Key)

	if err != nil {
		switch v := err.(type) {
		case *KeyNotFound:
			res.NotFound = v
		case *InvalidKey:
			res.InvalidKey = v
		default:
			return false, nil, err
		}
	} else {
		res.Success = &r
	}

	return err == nil, &res, nil
}

func (s *tchanKeyValueServer) handleHi(ctx thrift.Context, protocol athrift.TProtocol) (bool, athrift.TStruct, error) {
	var req KeyValueHiArgs
	var res KeyValueHiResult

	if err := req.Read(protocol); err != nil {
		return false, nil, err
	}

	err :=
		s.handler.Hi(ctx)

	if err != nil {
		return false, nil, err
	} else {
	}

	return err == nil, &res, nil
}

func (s *tchanKeyValueServer) handleSet(ctx thrift.Context, protocol athrift.TProtocol) (bool, athrift.TStruct, error) {
	var req KeyValueSetArgs
	var res KeyValueSetResult

	if err := req.Read(protocol); err != nil {
		return false, nil, err
	}

	err :=
		s.handler.Set(ctx, req.Key, req.Value)

	if err != nil {
		switch v := err.(type) {
		case *InvalidKey:
			res.InvalidKey = v
		default:
			return false, nil, err
		}
	} else {
	}

	return err == nil, &res, nil
}

type tchanBaseServiceClient struct {
	thriftService string
	client        thrift.TChanStreamingClient
}

func newTChanBaseServiceClient(thriftService string, client thrift.TChanStreamingClient) *tchanBaseServiceClient {
	return &tchanBaseServiceClient{
		thriftService,
		client,
	}
}

func NewTChanBaseServiceClient(client thrift.TChanStreamingClient) TChanBaseServiceClient {
	return newTChanBaseServiceClient("baseService", client)
}

func (c *tchanBaseServiceClient) HealthCheck(ctx thrift.Context) (string, error) {
	var resp BaseServiceHealthCheckResult
	args := BaseServiceHealthCheckArgs{}
	success, err := c.client.Call(ctx, c.thriftService, "HealthCheck", &args, &resp)
	if err == nil && !success {
	}

	return resp.GetSuccess(), err
}

type tchanBaseServiceServer struct {
	handler TChanBaseServiceServer
	common  thrift.TCommon
}

func newTChanBaseServiceServer(handler TChanBaseServiceServer) *tchanBaseServiceServer {
	return &tchanBaseServiceServer{
		handler,
		nil, /* common */
	}
}

func NewTChanBaseServiceServer(handler TChanBaseServiceServer) thrift.TChanStreamingServer {
	return newTChanBaseServiceServer(handler)
}

func (s *tchanBaseServiceServer) Service() string {
	return "baseService"
}

func (s *tchanBaseServiceServer) SetCommon(common thrift.TCommon) {
	s.common = common
}

func (s *tchanBaseServiceServer) Methods() []string {
	return []string{
		"HealthCheck",
	}
}

func (s *tchanBaseServiceServer) StreamingMethods() []string {
	return []string{}
}

func (s *tchanBaseServiceServer) HandleStreaming(ctx thrift.Context, call *tchannel.InboundCall) error {
	methodName := string(call.Operation())
	return fmt.Errorf("method %v not found in service %v", methodName, s.Service())
}

func (s *tchanBaseServiceServer) Handle(ctx thrift.Context, methodName string, protocol athrift.TProtocol) (bool, athrift.TStruct, error) {
	switch methodName {
	case "HealthCheck":
		return s.handleHealthCheck(ctx, protocol)

	default:
		return false, nil, fmt.Errorf("method %v not found in service %v", methodName, s.Service())
	}
}

func (s *tchanBaseServiceServer) handleHealthCheck(ctx thrift.Context, protocol athrift.TProtocol) (bool, athrift.TStruct, error) {
	var req BaseServiceHealthCheckArgs
	var res BaseServiceHealthCheckResult

	if err := req.Read(protocol); err != nil {
		return false, nil, err
	}

	r, err :=
		s.handler.HealthCheck(ctx)

	if err != nil {
		return false, nil, err
	} else {
		res.Success = &r
	}

	return err == nil, &res, nil
}
