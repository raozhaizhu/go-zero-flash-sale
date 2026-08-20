package logic

import (
	"context"

	"go-zero-flash-sale/app/usercenter/cmd/rpc/internal/svc"
	"go-zero-flash-sale/app/usercenter/cmd/rpc/usercenter"
	models "go-zero-flash-sale/app/usercenter/model"
	"go-zero-flash-sale/pkg/apperror"
	"go-zero-flash-sale/pkg/constants"
	"go-zero-flash-sale/pkg/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *usercenter.LoginRequest) (*usercenter.LoginResponse, error) {
	// 根据 AuthType 选择登录方式
	var userId int64
	var err error
	switch in.AuthType {
	case models.UserAuthTypeSystem:
		userId, err = l.loginByMobile(in.AuthKey, in.Password)
	default:
		return nil, apperror.WrapSrvErrf(apperror.SwitchNeverErr, "登录逻辑存在错误")
	}

	// 调用内部逻辑生成 token
	generateTokenLogic := NewGenerateTokenLogic(l.ctx, l.svcCtx)
	tokenResp, err := generateTokenLogic.GenerateToken(&usercenter.GenerateTokenRequest{
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}

	return &usercenter.LoginResponse{
		AccessToken:  tokenResp.AccessToken,
		AccessExpire: tokenResp.AccessExpire,
		RefreshAfter: tokenResp.RefreshAfter,
	}, nil
}

func (l *LoginLogic) loginByMobile(mobile, password string) (int64, error) {
	user, err := l.svcCtx.UserModel.FindOneByMobile(l.ctx, mobile)
	if err != nil {
		if err == models.ErrNotFound {
			return constants.InvalidUserId, apperror.ErrUserNotFound
		}
		return constants.InvalidUserId, apperror.WrapSrvErrf(err, "按手机登录时出现数据库错误")
	}

	err = util.CheckPassword(password, user.Password)
	if err != nil {
		return constants.InvalidUserId, apperror.WrapSrvErrf(err, "校验密码时出现内部错误")
	}

	return user.Id, nil
}
