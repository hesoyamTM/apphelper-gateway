package clients

import (
	"context"
	"fmt"

	"github.com/hesoyamTM/apphelper-sso/pkg/authorization"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func NewUIDInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		uid, ok := ctx.Value(authorization.Uid).(int64)
		if !ok {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		ctx = metadata.AppendToOutgoingContext(ctx, "uid", fmt.Sprintf("%d", uid))

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
