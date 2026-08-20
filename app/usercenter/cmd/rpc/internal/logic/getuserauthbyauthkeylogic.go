package logic

import (
	"context"

	"go-zero-flash-sale/app/usercenter/cmd/rpc/internal/svc"
	"go-zero-flash-sale/app/usercenter/cmd/rpc/usercenter"
	models "go-zero-flash-sale/app/usercenter/model"
	"go-zero-flash-sale/pkg/apperror"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserAuthByAuthKeyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserAuthByAuthKeyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserAuthByAuthKeyLogic {
	return &GetUserAuthByAuthKeyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserAuthByAuthKeyLogic) GetUserAuthByAuthKey(in *usercenter.GetUserAuthByAuthKeyRequest) (*usercenter.GetUserAuthByAuthKeyResponse, error) {
	userAuth, err := l.svcCtx.UserAuthModel.FindOneByAuthTypeAuthKey(l.ctx, in.AuthType, in.AuthKey)
	if err != nil {
		if err == models.ErrNotFound {
			return nil, apperror.ErrUserAuthNotFound
		}
		return nil, apperror.WrapSrvErrf(err, "查询 userAuth 时出现数据库错误, authType:%v authKey:%v", in.AuthType, in.AuthKey)
	}

	return &usercenter.GetUserAuthByAuthKeyResponse{
		UserAuth: &usercenter.UserAuth{
			Id:       userAuth.Id,
			UserId:   userAuth.UserId,
			AuthType: userAuth.AuthType,
			AuthKey:  userAuth.AuthKey,
		},
	}, nil
}
