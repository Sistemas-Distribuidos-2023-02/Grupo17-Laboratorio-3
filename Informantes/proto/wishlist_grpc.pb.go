// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package OMS

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// OMSClient is the client API for OMS service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OMSClient interface {
	NotifyBidirectional(ctx context.Context, opts ...grpc.CallOption) (OMS_NotifyBidirectionalClient, error)
}

type oMSClient struct {
	cc grpc.ClientConnInterface
}

func NewOMSClient(cc grpc.ClientConnInterface) OMSClient {
	return &oMSClient{cc}
}

func (c *oMSClient) NotifyBidirectional(ctx context.Context, opts ...grpc.CallOption) (OMS_NotifyBidirectionalClient, error) {
	stream, err := c.cc.NewStream(ctx, &OMS_ServiceDesc.Streams[0], "/OMS.OMS/NotifyBidirectional", opts...)
	if err != nil {
		return nil, err
	}
	x := &oMSNotifyBidirectionalClient{stream}
	return x, nil
}

type OMS_NotifyBidirectionalClient interface {
	Send(*Request) error
	Recv() (*Response, error)
	grpc.ClientStream
}

type oMSNotifyBidirectionalClient struct {
	grpc.ClientStream
}

func (x *oMSNotifyBidirectionalClient) Send(m *Request) error {
	return x.ClientStream.SendMsg(m)
}

func (x *oMSNotifyBidirectionalClient) Recv() (*Response, error) {
	m := new(Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// OMSServer is the server API for OMS service.
// All implementations must embed UnimplementedOMSServer
// for forward compatibility
type OMSServer interface {
	NotifyBidirectional(OMS_NotifyBidirectionalServer) error
	mustEmbedUnimplementedOMSServer()
}

// UnimplementedOMSServer must be embedded to have forward compatible implementations.
type UnimplementedOMSServer struct {
}

func (UnimplementedOMSServer) NotifyBidirectional(OMS_NotifyBidirectionalServer) error {
	return status.Errorf(codes.Unimplemented, "method NotifyBidirectional not implemented")
}
func (UnimplementedOMSServer) mustEmbedUnimplementedOMSServer() {}

// UnsafeOMSServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OMSServer will
// result in compilation errors.
type UnsafeOMSServer interface {
	mustEmbedUnimplementedOMSServer()
}

func RegisterOMSServer(s grpc.ServiceRegistrar, srv OMSServer) {
	s.RegisterService(&OMS_ServiceDesc, srv)
}

func _OMS_NotifyBidirectional_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(OMSServer).NotifyBidirectional(&oMSNotifyBidirectionalServer{stream})
}

type OMS_NotifyBidirectionalServer interface {
	Send(*Response) error
	Recv() (*Request, error)
	grpc.ServerStream
}

type oMSNotifyBidirectionalServer struct {
	grpc.ServerStream
}

func (x *oMSNotifyBidirectionalServer) Send(m *Response) error {
	return x.ServerStream.SendMsg(m)
}

func (x *oMSNotifyBidirectionalServer) Recv() (*Request, error) {
	m := new(Request)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// OMS_ServiceDesc is the grpc.ServiceDesc for OMS service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var OMS_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "OMS.OMS",
	HandlerType: (*OMSServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "NotifyBidirectional",
			Handler:       _OMS_NotifyBidirectional_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "wishlist.proto",
}
