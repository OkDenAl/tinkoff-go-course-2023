// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.6.1
// source: service.proto

package base

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	AdService_CreateAd_FullMethodName       = "/ad.AdService/CreateAd"
	AdService_ChangeAdStatus_FullMethodName = "/ad.AdService/ChangeAdStatus"
	AdService_UpdateAd_FullMethodName       = "/ad.AdService/UpdateAd"
	AdService_GetAdById_FullMethodName      = "/ad.AdService/GetAdById"
	AdService_GetAdByTitle_FullMethodName   = "/ad.AdService/GetAdByTitle"
	AdService_ListAds_FullMethodName        = "/ad.AdService/ListAds"
	AdService_DeleteAd_FullMethodName       = "/ad.AdService/DeleteAd"
)

// AdServiceClient is the client API for AdService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AdServiceClient interface {
	CreateAd(ctx context.Context, in *CreateAdRequest, opts ...grpc.CallOption) (*AdResponse, error)
	ChangeAdStatus(ctx context.Context, in *ChangeAdStatusRequest, opts ...grpc.CallOption) (*AdResponse, error)
	UpdateAd(ctx context.Context, in *UpdateAdRequest, opts ...grpc.CallOption) (*AdResponse, error)
	GetAdById(ctx context.Context, in *GetAdByIdRequest, opts ...grpc.CallOption) (*AdResponse, error)
	GetAdByTitle(ctx context.Context, in *GetAdByTitleRequest, opts ...grpc.CallOption) (*ListAdResponse, error)
	ListAds(ctx context.Context, in *Filters, opts ...grpc.CallOption) (*ListAdResponse, error)
	DeleteAd(ctx context.Context, in *DeleteAdRequest, opts ...grpc.CallOption) (*empty.Empty, error)
}

type adServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAdServiceClient(cc grpc.ClientConnInterface) AdServiceClient {
	return &adServiceClient{cc}
}

func (c *adServiceClient) CreateAd(ctx context.Context, in *CreateAdRequest, opts ...grpc.CallOption) (*AdResponse, error) {
	out := new(AdResponse)
	err := c.cc.Invoke(ctx, AdService_CreateAd_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adServiceClient) ChangeAdStatus(ctx context.Context, in *ChangeAdStatusRequest, opts ...grpc.CallOption) (*AdResponse, error) {
	out := new(AdResponse)
	err := c.cc.Invoke(ctx, AdService_ChangeAdStatus_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adServiceClient) UpdateAd(ctx context.Context, in *UpdateAdRequest, opts ...grpc.CallOption) (*AdResponse, error) {
	out := new(AdResponse)
	err := c.cc.Invoke(ctx, AdService_UpdateAd_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adServiceClient) GetAdById(ctx context.Context, in *GetAdByIdRequest, opts ...grpc.CallOption) (*AdResponse, error) {
	out := new(AdResponse)
	err := c.cc.Invoke(ctx, AdService_GetAdById_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adServiceClient) GetAdByTitle(ctx context.Context, in *GetAdByTitleRequest, opts ...grpc.CallOption) (*ListAdResponse, error) {
	out := new(ListAdResponse)
	err := c.cc.Invoke(ctx, AdService_GetAdByTitle_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adServiceClient) ListAds(ctx context.Context, in *Filters, opts ...grpc.CallOption) (*ListAdResponse, error) {
	out := new(ListAdResponse)
	err := c.cc.Invoke(ctx, AdService_ListAds_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adServiceClient) DeleteAd(ctx context.Context, in *DeleteAdRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, AdService_DeleteAd_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AdServiceServer is the server API for AdService service.
// All implementations must embed UnimplementedAdServiceServer
// for forward compatibility
type AdServiceServer interface {
	CreateAd(context.Context, *CreateAdRequest) (*AdResponse, error)
	ChangeAdStatus(context.Context, *ChangeAdStatusRequest) (*AdResponse, error)
	UpdateAd(context.Context, *UpdateAdRequest) (*AdResponse, error)
	GetAdById(context.Context, *GetAdByIdRequest) (*AdResponse, error)
	GetAdByTitle(context.Context, *GetAdByTitleRequest) (*ListAdResponse, error)
	ListAds(context.Context, *Filters) (*ListAdResponse, error)
	DeleteAd(context.Context, *DeleteAdRequest) (*empty.Empty, error)
	mustEmbedUnimplementedAdServiceServer()
}

// UnimplementedAdServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAdServiceServer struct {
}

func (UnimplementedAdServiceServer) CreateAd(context.Context, *CreateAdRequest) (*AdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAd not implemented")
}
func (UnimplementedAdServiceServer) ChangeAdStatus(context.Context, *ChangeAdStatusRequest) (*AdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangeAdStatus not implemented")
}
func (UnimplementedAdServiceServer) UpdateAd(context.Context, *UpdateAdRequest) (*AdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateAd not implemented")
}
func (UnimplementedAdServiceServer) GetAdById(context.Context, *GetAdByIdRequest) (*AdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAdById not implemented")
}
func (UnimplementedAdServiceServer) GetAdByTitle(context.Context, *GetAdByTitleRequest) (*ListAdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAdByTitle not implemented")
}
func (UnimplementedAdServiceServer) ListAds(context.Context, *Filters) (*ListAdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAds not implemented")
}
func (UnimplementedAdServiceServer) DeleteAd(context.Context, *DeleteAdRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAd not implemented")
}
func (UnimplementedAdServiceServer) mustEmbedUnimplementedAdServiceServer() {}

// UnsafeAdServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AdServiceServer will
// result in compilation errors.
type UnsafeAdServiceServer interface {
	mustEmbedUnimplementedAdServiceServer()
}

func RegisterAdServiceServer(s grpc.ServiceRegistrar, srv AdServiceServer) {
	s.RegisterService(&AdService_ServiceDesc, srv)
}

func _AdService_CreateAd_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateAdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdServiceServer).CreateAd(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdService_CreateAd_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdServiceServer).CreateAd(ctx, req.(*CreateAdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdService_ChangeAdStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChangeAdStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdServiceServer).ChangeAdStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdService_ChangeAdStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdServiceServer).ChangeAdStatus(ctx, req.(*ChangeAdStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdService_UpdateAd_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateAdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdServiceServer).UpdateAd(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdService_UpdateAd_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdServiceServer).UpdateAd(ctx, req.(*UpdateAdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdService_GetAdById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAdByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdServiceServer).GetAdById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdService_GetAdById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdServiceServer).GetAdById(ctx, req.(*GetAdByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdService_GetAdByTitle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAdByTitleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdServiceServer).GetAdByTitle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdService_GetAdByTitle_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdServiceServer).GetAdByTitle(ctx, req.(*GetAdByTitleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdService_ListAds_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Filters)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdServiceServer).ListAds(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdService_ListAds_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdServiceServer).ListAds(ctx, req.(*Filters))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdService_DeleteAd_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteAdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdServiceServer).DeleteAd(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdService_DeleteAd_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdServiceServer).DeleteAd(ctx, req.(*DeleteAdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AdService_ServiceDesc is the grpc.ServiceDesc for AdService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AdService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ad.AdService",
	HandlerType: (*AdServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateAd",
			Handler:    _AdService_CreateAd_Handler,
		},
		{
			MethodName: "ChangeAdStatus",
			Handler:    _AdService_ChangeAdStatus_Handler,
		},
		{
			MethodName: "UpdateAd",
			Handler:    _AdService_UpdateAd_Handler,
		},
		{
			MethodName: "GetAdById",
			Handler:    _AdService_GetAdById_Handler,
		},
		{
			MethodName: "GetAdByTitle",
			Handler:    _AdService_GetAdByTitle_Handler,
		},
		{
			MethodName: "ListAds",
			Handler:    _AdService_ListAds_Handler,
		},
		{
			MethodName: "DeleteAd",
			Handler:    _AdService_DeleteAd_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service.proto",
}

const (
	UserService_CreateUser_FullMethodName     = "/ad.UserService/CreateUser"
	UserService_ChangeNickname_FullMethodName = "/ad.UserService/ChangeNickname"
	UserService_GetUser_FullMethodName        = "/ad.UserService/GetUser"
	UserService_DeleteUser_FullMethodName     = "/ad.UserService/DeleteUser"
)

// UserServiceClient is the client API for UserService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserServiceClient interface {
	CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*UserResponse, error)
	ChangeNickname(ctx context.Context, in *ChangeNicknameRequest, opts ...grpc.CallOption) (*UserResponse, error)
	GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*UserResponse, error)
	DeleteUser(ctx context.Context, in *DeleteUserRequest, opts ...grpc.CallOption) (*empty.Empty, error)
}

type userServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUserServiceClient(cc grpc.ClientConnInterface) UserServiceClient {
	return &userServiceClient{cc}
}

func (c *userServiceClient) CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*UserResponse, error) {
	out := new(UserResponse)
	err := c.cc.Invoke(ctx, UserService_CreateUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) ChangeNickname(ctx context.Context, in *ChangeNicknameRequest, opts ...grpc.CallOption) (*UserResponse, error) {
	out := new(UserResponse)
	err := c.cc.Invoke(ctx, UserService_ChangeNickname_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*UserResponse, error) {
	out := new(UserResponse)
	err := c.cc.Invoke(ctx, UserService_GetUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) DeleteUser(ctx context.Context, in *DeleteUserRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, UserService_DeleteUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserServiceServer is the server API for UserService service.
// All implementations must embed UnimplementedUserServiceServer
// for forward compatibility
type UserServiceServer interface {
	CreateUser(context.Context, *CreateUserRequest) (*UserResponse, error)
	ChangeNickname(context.Context, *ChangeNicknameRequest) (*UserResponse, error)
	GetUser(context.Context, *GetUserRequest) (*UserResponse, error)
	DeleteUser(context.Context, *DeleteUserRequest) (*empty.Empty, error)
	mustEmbedUnimplementedUserServiceServer()
}

// UnimplementedUserServiceServer must be embedded to have forward compatible implementations.
type UnimplementedUserServiceServer struct {
}

func (UnimplementedUserServiceServer) CreateUser(context.Context, *CreateUserRequest) (*UserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}
func (UnimplementedUserServiceServer) ChangeNickname(context.Context, *ChangeNicknameRequest) (*UserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangeNickname not implemented")
}
func (UnimplementedUserServiceServer) GetUser(context.Context, *GetUserRequest) (*UserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUser not implemented")
}
func (UnimplementedUserServiceServer) DeleteUser(context.Context, *DeleteUserRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUser not implemented")
}
func (UnimplementedUserServiceServer) mustEmbedUnimplementedUserServiceServer() {}

// UnsafeUserServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserServiceServer will
// result in compilation errors.
type UnsafeUserServiceServer interface {
	mustEmbedUnimplementedUserServiceServer()
}

func RegisterUserServiceServer(s grpc.ServiceRegistrar, srv UserServiceServer) {
	s.RegisterService(&UserService_ServiceDesc, srv)
}

func _UserService_CreateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_CreateUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).CreateUser(ctx, req.(*CreateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_ChangeNickname_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChangeNicknameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).ChangeNickname(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_ChangeNickname_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).ChangeNickname(ctx, req.(*ChangeNicknameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_GetUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetUser(ctx, req.(*GetUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_DeleteUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).DeleteUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_DeleteUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).DeleteUser(ctx, req.(*DeleteUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UserService_ServiceDesc is the grpc.ServiceDesc for UserService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ad.UserService",
	HandlerType: (*UserServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateUser",
			Handler:    _UserService_CreateUser_Handler,
		},
		{
			MethodName: "ChangeNickname",
			Handler:    _UserService_ChangeNickname_Handler,
		},
		{
			MethodName: "GetUser",
			Handler:    _UserService_GetUser_Handler,
		},
		{
			MethodName: "DeleteUser",
			Handler:    _UserService_DeleteUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service.proto",
}