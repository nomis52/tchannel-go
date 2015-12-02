// @generated Code generated by thrift-gen. Do not modify.

// Package test is generated code used to make or handle TChannel calls using Thrift.
package test

import (
	"fmt"

	athrift "github.com/apache/thrift/lib/go/thrift"
	"github.com/uber/tchannel-go/thrift"
)

// Interfaces for the service and client for the services defined in the IDL.

// TChanBase is the interface that defines the server handler and client interface.
type TChanBase interface {
	BaseCall(ctx thrift.Context) error
}

// TChanFirst is the interface that defines the server handler and client interface.
type TChanFirst interface {
	TChanBase

	AppError(ctx thrift.Context) error
	Echo(ctx thrift.Context, msg string) (string, error)
	Healthcheck(ctx thrift.Context) (*HealthCheckRes, error)
}

// TChanSecond is the interface that defines the server handler and client interface.
type TChanSecond interface {
	Test(ctx thrift.Context) error
}

// Implementation of a client and service handler.

type tchanBaseClient struct {
	thriftService string
	client        thrift.TChanClient
}

func newTChanBaseClient(thriftService string, client thrift.TChanClient) *tchanBaseClient {
	return &tchanBaseClient{
		thriftService,
		client,
	}
}

// NewTChanBaseClient creates a client that can be used to make remote calls.
func NewTChanBaseClient(client thrift.TChanClient) TChanBase {
	return newTChanBaseClient("Base", client)
}

func (c *tchanBaseClient) BaseCall(ctx thrift.Context) error {
	var resp BaseBaseCallResult
	args := BaseBaseCallArgs{}
	success, err := c.client.Call(ctx, c.thriftService, "BaseCall", &args, &resp)
	if err == nil && !success {
	}

	return err
}

type tchanBaseServer struct {
	handler TChanBase
}

func newTChanBaseServer(handler TChanBase) *tchanBaseServer {
	return &tchanBaseServer{
		handler,
	}
}

// NewTChanBaseServer wraps a handler for TChanBase so it can be
// registered with a thrift.Server.
func NewTChanBaseServer(handler TChanBase) thrift.TChanServer {
	return newTChanBaseServer(handler)
}

func (s *tchanBaseServer) Service() string {
	return "Base"
}

func (s *tchanBaseServer) Methods() []string {
	return []string{
		"BaseCall",
	}
}

func (s *tchanBaseServer) Handle(ctx thrift.Context, methodName string, protocol athrift.TProtocol) (bool, athrift.TStruct, error) {
	args, err := s.GetArgs(methodName, protocol)
	if err != nil {
		return false, nil, err
	}
	return s.HandleArgs(ctx, methodName, args)
}

func (s *tchanBaseServer) GetArgs(methodName string, protocol athrift.TProtocol) (args interface{}, err error) {
	switch methodName {
	case "BaseCall":
		args, err = s.readBaseCall(protocol)
	default:
		err = fmt.Errorf("method %v not found in service %v", methodName, s.Service())
	}
	return
}

func (s *tchanBaseServer) HandleArgs(ctx thrift.Context, methodName string, args interface{}) (bool, athrift.TStruct, error) {
	switch methodName {
	case "BaseCall":
		return s.handleBaseCall(ctx, args.(BaseBaseCallArgs))
	default:
		return false, nil, fmt.Errorf("method %v not found in service %v", methodName, s.Service())
	}
}

func (s *tchanBaseServer) readBaseCall(protocol athrift.TProtocol) (interface{}, error) {
	var req BaseBaseCallArgs

	if err := req.Read(protocol); err != nil {
		return nil, err
	}
	return req, nil
}

func (s *tchanBaseServer) handleBaseCall(ctx thrift.Context, req BaseBaseCallArgs) (bool, athrift.TStruct, error) {
	var res BaseBaseCallResult
	err :=
		s.handler.BaseCall(ctx)

	if err != nil {
		return false, nil, err
	} else {
	}

	return err == nil, &res, nil
}

type tchanFirstClient struct {
	tchanBaseClient

	thriftService string
	client        thrift.TChanClient
}

func newTChanFirstClient(thriftService string, client thrift.TChanClient) *tchanFirstClient {
	return &tchanFirstClient{
		*newTChanBaseClient(thriftService, client),
		thriftService,
		client,
	}
}

// NewTChanFirstClient creates a client that can be used to make remote calls.
func NewTChanFirstClient(client thrift.TChanClient) TChanFirst {
	return newTChanFirstClient("First", client)
}

func (c *tchanFirstClient) AppError(ctx thrift.Context) error {
	var resp FirstAppErrorResult
	args := FirstAppErrorArgs{}
	success, err := c.client.Call(ctx, c.thriftService, "AppError", &args, &resp)
	if err == nil && !success {
	}

	return err
}

func (c *tchanFirstClient) Echo(ctx thrift.Context, msg string) (string, error) {
	var resp FirstEchoResult
	args := FirstEchoArgs{
		Msg: msg,
	}
	success, err := c.client.Call(ctx, c.thriftService, "Echo", &args, &resp)
	if err == nil && !success {
	}

	return resp.GetSuccess(), err
}

func (c *tchanFirstClient) Healthcheck(ctx thrift.Context) (*HealthCheckRes, error) {
	var resp FirstHealthcheckResult
	args := FirstHealthcheckArgs{}
	success, err := c.client.Call(ctx, c.thriftService, "Healthcheck", &args, &resp)
	if err == nil && !success {
	}

	return resp.GetSuccess(), err
}

type tchanFirstServer struct {
	tchanBaseServer

	handler TChanFirst
}

func newTChanFirstServer(handler TChanFirst) *tchanFirstServer {
	return &tchanFirstServer{
		*newTChanBaseServer(handler),
		handler,
	}
}

// NewTChanFirstServer wraps a handler for TChanFirst so it can be
// registered with a thrift.Server.
func NewTChanFirstServer(handler TChanFirst) thrift.TChanServer {
	return newTChanFirstServer(handler)
}

func (s *tchanFirstServer) Service() string {
	return "First"
}

func (s *tchanFirstServer) Methods() []string {
	return []string{
		"AppError",
		"Echo",
		"Healthcheck",

		"BaseCall",
	}
}

func (s *tchanFirstServer) Handle(ctx thrift.Context, methodName string, protocol athrift.TProtocol) (bool, athrift.TStruct, error) {
	args, err := s.GetArgs(methodName, protocol)
	if err != nil {
		return false, nil, err
	}
	return s.HandleArgs(ctx, methodName, args)
}

func (s *tchanFirstServer) GetArgs(methodName string, protocol athrift.TProtocol) (args interface{}, err error) {
	switch methodName {
	case "AppError":
		args, err = s.readAppError(protocol)
	case "Echo":
		args, err = s.readEcho(protocol)
	case "Healthcheck":
		args, err = s.readHealthcheck(protocol)
	case "BaseCall":
		args, err = s.readBaseCall(protocol)
	default:
		err = fmt.Errorf("method %v not found in service %v", methodName, s.Service())
	}
	return
}

func (s *tchanFirstServer) HandleArgs(ctx thrift.Context, methodName string, args interface{}) (bool, athrift.TStruct, error) {
	switch methodName {
	case "AppError":
		return s.handleAppError(ctx, args.(FirstAppErrorArgs))
	case "Echo":
		return s.handleEcho(ctx, args.(FirstEchoArgs))
	case "Healthcheck":
		return s.handleHealthcheck(ctx, args.(FirstHealthcheckArgs))
	case "BaseCall":
		return s.handleBaseCall(ctx, args.(BaseBaseCallArgs))
	default:
		return false, nil, fmt.Errorf("method %v not found in service %v", methodName, s.Service())
	}
}

func (s *tchanFirstServer) readAppError(protocol athrift.TProtocol) (interface{}, error) {
	var req FirstAppErrorArgs

	if err := req.Read(protocol); err != nil {
		return nil, err
	}
	return req, nil
}

func (s *tchanFirstServer) handleAppError(ctx thrift.Context, req FirstAppErrorArgs) (bool, athrift.TStruct, error) {
	var res FirstAppErrorResult
	err :=
		s.handler.AppError(ctx)

	if err != nil {
		return false, nil, err
	} else {
	}

	return err == nil, &res, nil
}

func (s *tchanFirstServer) readEcho(protocol athrift.TProtocol) (interface{}, error) {
	var req FirstEchoArgs

	if err := req.Read(protocol); err != nil {
		return nil, err
	}
	return req, nil
}

func (s *tchanFirstServer) handleEcho(ctx thrift.Context, req FirstEchoArgs) (bool, athrift.TStruct, error) {
	var res FirstEchoResult
	r, err :=
		s.handler.Echo(ctx, req.Msg)

	if err != nil {
		return false, nil, err
	} else {
		res.Success = &r
	}

	return err == nil, &res, nil
}

func (s *tchanFirstServer) readHealthcheck(protocol athrift.TProtocol) (interface{}, error) {
	var req FirstHealthcheckArgs

	if err := req.Read(protocol); err != nil {
		return nil, err
	}
	return req, nil
}

func (s *tchanFirstServer) handleHealthcheck(ctx thrift.Context, req FirstHealthcheckArgs) (bool, athrift.TStruct, error) {
	var res FirstHealthcheckResult
	r, err :=
		s.handler.Healthcheck(ctx)

	if err != nil {
		return false, nil, err
	} else {
		res.Success = r
	}

	return err == nil, &res, nil
}

type tchanSecondClient struct {
	thriftService string
	client        thrift.TChanClient
}

func newTChanSecondClient(thriftService string, client thrift.TChanClient) *tchanSecondClient {
	return &tchanSecondClient{
		thriftService,
		client,
	}
}

// NewTChanSecondClient creates a client that can be used to make remote calls.
func NewTChanSecondClient(client thrift.TChanClient) TChanSecond {
	return newTChanSecondClient("Second", client)
}

func (c *tchanSecondClient) Test(ctx thrift.Context) error {
	var resp SecondTestResult
	args := SecondTestArgs{}
	success, err := c.client.Call(ctx, c.thriftService, "Test", &args, &resp)
	if err == nil && !success {
	}

	return err
}

type tchanSecondServer struct {
	handler TChanSecond
}

func newTChanSecondServer(handler TChanSecond) *tchanSecondServer {
	return &tchanSecondServer{
		handler,
	}
}

// NewTChanSecondServer wraps a handler for TChanSecond so it can be
// registered with a thrift.Server.
func NewTChanSecondServer(handler TChanSecond) thrift.TChanServer {
	return newTChanSecondServer(handler)
}

func (s *tchanSecondServer) Service() string {
	return "Second"
}

func (s *tchanSecondServer) Methods() []string {
	return []string{
		"Test",
	}
}

func (s *tchanSecondServer) Handle(ctx thrift.Context, methodName string, protocol athrift.TProtocol) (bool, athrift.TStruct, error) {
	args, err := s.GetArgs(methodName, protocol)
	if err != nil {
		return false, nil, err
	}
	return s.HandleArgs(ctx, methodName, args)
}

func (s *tchanSecondServer) GetArgs(methodName string, protocol athrift.TProtocol) (args interface{}, err error) {
	switch methodName {
	case "Test":
		args, err = s.readTest(protocol)
	default:
		err = fmt.Errorf("method %v not found in service %v", methodName, s.Service())
	}
	return
}

func (s *tchanSecondServer) HandleArgs(ctx thrift.Context, methodName string, args interface{}) (bool, athrift.TStruct, error) {
	switch methodName {
	case "Test":
		return s.handleTest(ctx, args.(SecondTestArgs))
	default:
		return false, nil, fmt.Errorf("method %v not found in service %v", methodName, s.Service())
	}
}

func (s *tchanSecondServer) readTest(protocol athrift.TProtocol) (interface{}, error) {
	var req SecondTestArgs

	if err := req.Read(protocol); err != nil {
		return nil, err
	}
	return req, nil
}

func (s *tchanSecondServer) handleTest(ctx thrift.Context, req SecondTestArgs) (bool, athrift.TStruct, error) {
	var res SecondTestResult
	err :=
		s.handler.Test(ctx)

	if err != nil {
		return false, nil, err
	} else {
	}

	return err == nil, &res, nil
}
