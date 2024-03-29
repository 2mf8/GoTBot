// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.1
// source: http.proto

package dto

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

const (
	HttpService_Login_FullMethodName                      = "/dto.HttpService/Login"
	HttpService_Logout_FullMethodName                     = "/dto.HttpService/Logout"
	HttpService_GetShopItemAll_FullMethodName             = "/dto.HttpService/GetShopItemAll"
	HttpService_DeleteShopItem_FullMethodName             = "/dto.HttpService/DeleteShopItem"
	HttpService_AddAndUpdateShopItemByItem_FullMethodName = "/dto.HttpService/AddAndUpdateShopItemByItem"
)

// HttpServiceClient is the client API for HttpService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type HttpServiceClient interface {
	Login(ctx context.Context, in *CodeLoginReq, opts ...grpc.CallOption) (*CodeLoginResp, error)
	Logout(ctx context.Context, in *LogoutReq, opts ...grpc.CallOption) (*LogoutResp, error)
	GetShopItemAll(ctx context.Context, in *GetShopItemAllReq, opts ...grpc.CallOption) (*GetShopItemAllResp, error)
	DeleteShopItem(ctx context.Context, in *DeleteShopItemReq, opts ...grpc.CallOption) (*DeleteShopItemResp, error)
	AddAndUpdateShopItemByItem(ctx context.Context, in *AddAndUpdateShopItemByItemReq, opts ...grpc.CallOption) (*AddAndUpdateShopItemByItemResp, error)
}

type httpServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewHttpServiceClient(cc grpc.ClientConnInterface) HttpServiceClient {
	return &httpServiceClient{cc}
}

func (c *httpServiceClient) Login(ctx context.Context, in *CodeLoginReq, opts ...grpc.CallOption) (*CodeLoginResp, error) {
	out := new(CodeLoginResp)
	err := c.cc.Invoke(ctx, HttpService_Login_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *httpServiceClient) Logout(ctx context.Context, in *LogoutReq, opts ...grpc.CallOption) (*LogoutResp, error) {
	out := new(LogoutResp)
	err := c.cc.Invoke(ctx, HttpService_Logout_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *httpServiceClient) GetShopItemAll(ctx context.Context, in *GetShopItemAllReq, opts ...grpc.CallOption) (*GetShopItemAllResp, error) {
	out := new(GetShopItemAllResp)
	err := c.cc.Invoke(ctx, HttpService_GetShopItemAll_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *httpServiceClient) DeleteShopItem(ctx context.Context, in *DeleteShopItemReq, opts ...grpc.CallOption) (*DeleteShopItemResp, error) {
	out := new(DeleteShopItemResp)
	err := c.cc.Invoke(ctx, HttpService_DeleteShopItem_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *httpServiceClient) AddAndUpdateShopItemByItem(ctx context.Context, in *AddAndUpdateShopItemByItemReq, opts ...grpc.CallOption) (*AddAndUpdateShopItemByItemResp, error) {
	out := new(AddAndUpdateShopItemByItemResp)
	err := c.cc.Invoke(ctx, HttpService_AddAndUpdateShopItemByItem_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HttpServiceServer is the server API for HttpService service.
// All implementations must embed UnimplementedHttpServiceServer
// for forward compatibility
type HttpServiceServer interface {
	Login(context.Context, *CodeLoginReq) (*CodeLoginResp, error)
	Logout(context.Context, *LogoutReq) (*LogoutResp, error)
	GetShopItemAll(context.Context, *GetShopItemAllReq) (*GetShopItemAllResp, error)
	DeleteShopItem(context.Context, *DeleteShopItemReq) (*DeleteShopItemResp, error)
	AddAndUpdateShopItemByItem(context.Context, *AddAndUpdateShopItemByItemReq) (*AddAndUpdateShopItemByItemResp, error)
	mustEmbedUnimplementedHttpServiceServer()
}

// UnimplementedHttpServiceServer must be embedded to have forward compatible implementations.
type UnimplementedHttpServiceServer struct {
}

func (UnimplementedHttpServiceServer) Login(context.Context, *CodeLoginReq) (*CodeLoginResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedHttpServiceServer) Logout(context.Context, *LogoutReq) (*LogoutResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Logout not implemented")
}
func (UnimplementedHttpServiceServer) GetShopItemAll(context.Context, *GetShopItemAllReq) (*GetShopItemAllResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetShopItemAll not implemented")
}
func (UnimplementedHttpServiceServer) DeleteShopItem(context.Context, *DeleteShopItemReq) (*DeleteShopItemResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteShopItem not implemented")
}
func (UnimplementedHttpServiceServer) AddAndUpdateShopItemByItem(context.Context, *AddAndUpdateShopItemByItemReq) (*AddAndUpdateShopItemByItemResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddAndUpdateShopItemByItem not implemented")
}
func (UnimplementedHttpServiceServer) mustEmbedUnimplementedHttpServiceServer() {}

// UnsafeHttpServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to HttpServiceServer will
// result in compilation errors.
type UnsafeHttpServiceServer interface {
	mustEmbedUnimplementedHttpServiceServer()
}

func RegisterHttpServiceServer(s grpc.ServiceRegistrar, srv HttpServiceServer) {
	s.RegisterService(&HttpService_ServiceDesc, srv)
}

func _HttpService_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CodeLoginReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HttpServiceServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: HttpService_Login_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HttpServiceServer).Login(ctx, req.(*CodeLoginReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _HttpService_Logout_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogoutReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HttpServiceServer).Logout(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: HttpService_Logout_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HttpServiceServer).Logout(ctx, req.(*LogoutReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _HttpService_GetShopItemAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetShopItemAllReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HttpServiceServer).GetShopItemAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: HttpService_GetShopItemAll_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HttpServiceServer).GetShopItemAll(ctx, req.(*GetShopItemAllReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _HttpService_DeleteShopItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteShopItemReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HttpServiceServer).DeleteShopItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: HttpService_DeleteShopItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HttpServiceServer).DeleteShopItem(ctx, req.(*DeleteShopItemReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _HttpService_AddAndUpdateShopItemByItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddAndUpdateShopItemByItemReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HttpServiceServer).AddAndUpdateShopItemByItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: HttpService_AddAndUpdateShopItemByItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HttpServiceServer).AddAndUpdateShopItemByItem(ctx, req.(*AddAndUpdateShopItemByItemReq))
	}
	return interceptor(ctx, in, info, handler)
}

// HttpService_ServiceDesc is the grpc.ServiceDesc for HttpService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var HttpService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "dto.HttpService",
	HandlerType: (*HttpServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Login",
			Handler:    _HttpService_Login_Handler,
		},
		{
			MethodName: "Logout",
			Handler:    _HttpService_Logout_Handler,
		},
		{
			MethodName: "GetShopItemAll",
			Handler:    _HttpService_GetShopItemAll_Handler,
		},
		{
			MethodName: "DeleteShopItem",
			Handler:    _HttpService_DeleteShopItem_Handler,
		},
		{
			MethodName: "AddAndUpdateShopItemByItem",
			Handler:    _HttpService_AddAndUpdateShopItemByItem_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "http.proto",
}
