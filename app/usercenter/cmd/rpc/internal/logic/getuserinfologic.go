package logic

import (
	"context"

	"go-zero-flash-sale/app/usercenter/cmd/rpc/internal/svc"
	"go-zero-flash-sale/app/usercenter/cmd/rpc/usercenter"
	models "go-zero-flash-sale/app/usercenter/model"
	"go-zero-flash-sale/pkg/apperror"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoLogic) GetUserInfo(in *usercenter.GetUserInfoRequest) (*usercenter.GetUserInfoResponse, error) {
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == models.ErrNotFound {
			return nil, apperror.ErrUserNotFound
		}
		return nil, apperror.WrapSrvErrf(err, "查询 user 时出现数据库错误, userId:%v", in.Id)
	}

	return &usercenter.GetUserInfoResponse{
		User: &usercenter.User{
			Id:       user.Id,
			Mobile:   user.Mobile,
			Nickname: user.Nickname,
			Sex:      user.Sex,
			Avatar:   user.Avatar,
			Info:     user.Info,
		},
	}, nil
}
