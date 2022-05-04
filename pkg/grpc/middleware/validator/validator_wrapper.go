// Copyright 2022 Monoskope Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package validator

import (
	"context"
	"github.com/finleap-connect/monoskope/pkg/api/eventsourcing/commands"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type validator interface {
	Validate() error
}

// UnaryServerInterceptor returns a new unary server interceptor that validates incoming messages.
//
// Invalid messages will be rejected with `InvalidArgument` before reaching any userspace handlers.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if err := validate(req); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		return handler(ctx, req)
	}
}

// UnaryClientInterceptor returns a new unary client interceptor that validates outgoing messages.
//
// Invalid messages will be rejected with `InvalidArgument` before sending the request to server.
func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if err := validate(req); err != nil {
			return status.Errorf(codes.InvalidArgument, err.Error())
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// StreamServerInterceptor returns a new streaming server interceptor that validates incoming messages.
//
// The stage at which invalid messages will be rejected with `InvalidArgument` varies based on the
// type of the RPC. For `ServerStream` (1:m) requests, it will happen before reaching any userspace
// handlers. For `ClientStream` (n:1) or `BidiStream` (n:m) RPCs, the messages will be rejected on
// calls to `stream.Recv()`.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		wrapper := &recvWrapper{stream}
		return handler(srv, wrapper)
	}
}

type recvWrapper struct {
	grpc.ServerStream
}

func (s *recvWrapper) RecvMsg(m interface{}) error {
	if err := s.ServerStream.RecvMsg(m); err != nil {
		return err
	}
	if err := validate(m); err != nil {
		return status.Errorf(codes.InvalidArgument, err.Error())
	}
	return nil
}

func validate(req interface{}) error {
	if v, ok := req.(validator); ok {
		if err := v.Validate(); err != nil {
			return err
		}
	}

	// unpack command data for validation
	if c, ok := req.(*commands.Command); ok {
		if d, err := c.Data.UnmarshalNew(); err == nil {
			if v, ok := d.ProtoReflect().Interface().(validator); ok {
				if err := v.Validate(); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
