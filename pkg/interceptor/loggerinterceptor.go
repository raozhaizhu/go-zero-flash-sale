package interceptor

import (
	"context"
	"go-zero-flash-sale/pkg/apperror"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func LoggerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	// 处理请求
	resp, err = handler(ctx, req)

	if err != nil {
		var bizErr *apperror.BizError

		if errors.As(err, &bizErr) {
			if bizErr.Code == apperror.CodeServerErr {
				// 如果是服务器内部错误，把信封里的底层错误（bizErr.err）打印出来
				logx.WithContext(ctx).Errorf("【RPC-SYS-ERR】 Code: %d, Msg: %s, Cause: %+v", bizErr.Code, bizErr.Msg, bizErr.Unwrap())
				err = status.Error(codes.Internal, "服务器繁忙,请稍后再试")
			} else {
				// 普通业务错误
				logx.WithContext(ctx).Infof("【RPC-BIZ-WARN】 业务拦截: %v", bizErr.Error())
				err = status.Error(codes.Code(bizErr.Code), bizErr.Msg)
			}
		} else {
			// 未知错误兜底
			logx.WithContext(ctx).Errorf("【RPC-UNKNOWN-ERR】 %+v", err)
			err = status.Error(codes.Internal, "服务器繁忙,请稍后再试")
		}
	}

	return resp, err
}
