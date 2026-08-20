package logic

import (
	"context"

	"go-zero-flash-sale/app/usercenter/cmd/rpc/internal/svc"
	"go-zero-flash-sale/app/usercenter/cmd/rpc/usercenter"
	models "go-zero-flash-sale/app/usercenter/model"
	"go-zero-flash-sale/pkg/apperror"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserAuthByUserIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserAuthByUserIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserAuthByUserIdLogic {
	return &GetUserAuthByUserIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserAuthByUserIdLogic) GetUserAuthByUserId(in *usercenter.GetUserAuthByUserIdRequest) (*usercenter.GetUserAuthByUserIdResponse, error) {
	userAuth, err := l.svcCtx.UserAuthModel.FindOneByUserIdAuthType(l.ctx, in.UserId, in.AuthType)
	if err != nil {
		if err == models.ErrNotFound {
			return nil, apperror.ErrUserAuthNotFound
		}
		return nil, apperror.WrapSrvErrf(err, "查询 userAuth 时出现数据库错误, userId:%v", in.UserId)
	}

	return &usercenter.GetUserAuthByUserIdResponse{
		UserAuth: &usercenter.UserAuth{
			Id:       userAuth.Id,
			UserId:   userAuth.UserId,
			AuthType: userAuth.AuthType,
			AuthKey:  userAuth.AuthKey,
		},
	}, nil
}
