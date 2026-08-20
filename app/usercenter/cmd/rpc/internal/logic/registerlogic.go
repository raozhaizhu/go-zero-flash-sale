package logic

import (
	"context"

	"go-zero-flash-sale/app/usercenter/cmd/rpc/internal/svc"
	"go-zero-flash-sale/app/usercenter/cmd/rpc/usercenter"
	"go-zero-flash-sale/pkg/apperror"
	"go-zero-flash-sale/pkg/util"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *usercenter.RegisterRequest) (*usercenter.RegisterResponse, error) {
	// 1. 业务逻辑校验：用户已存在
	_, err := l.svcCtx.UserModel.FindOneByMobile(l.ctx, in.Mobile)
	if err == nil {
		return nil, apperror.ErrUserAlreadyExits
	}

	// 2. 前置处理：昵称与密码哈希
	nickname := in.Nickname
	if len(nickname) == 0 {
		nickname = util.Krand(8, util.KindAll)
	}

	var hashedPassword string
	if len(in.Password) > 0 {
		hashedPassword, err = util.HashPassword(in.Password)
		if err != nil {
			return nil, apperror.WrapSrvErrf(err, "哈希密码失败, mobile: %s", in.Mobile)
		}
	}

	// 3. 调用 TxRunner 开启事务
	var userId int64
	err = l.svcCtx.TxRunner.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		userQuery := "INSERT INTO user (mobile, nickname, password) VALUES (?, ?, ?)"
		insertResult, err := session.ExecCtx(ctx, userQuery, in.Mobile, nickname, hashedPassword)
		if err != nil {
			return err
		}

		lastId, err := insertResult.LastInsertId()
		if err != nil {
			return err
		}
		userId = lastId

		authQuery := "INSERT INTO user_auth (user_id, auth_type, auth_key) VALUES (?, ?, ?)"
		_, err = session.ExecCtx(ctx, authQuery, userId, in.AuthType, in.AuthKey)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, apperror.WrapSrvErrf(err, "注册数据库事务执行失败, mobile: %s", in.Mobile)
	}

	// 4. 异步缓存预热
	asyncCtx := context.WithoutCancel(l.ctx)
	go func() {
		_, err = l.svcCtx.UserModel.FindOne(asyncCtx, userId)
		if err != nil {
			logx.WithContext(asyncCtx).Errorf("异步预热缓存失败, userId: %v, err: %v", userId, err)
		}
	}()

	// 5. 调用内部逻辑生成 Token
	generateTokenLogic := NewGenerateTokenLogic(l.ctx, l.svcCtx)
	tokenResp, err := generateTokenLogic.GenerateToken(&usercenter.GenerateTokenRequest{
		UserId: userId,
	})
	if err != nil {
		return nil, apperror.WrapSrvErrf(err, "GenerateToken 调用失败, userId: %d", userId)
	}

	return &usercenter.RegisterResponse{
		AccessToken:  tokenResp.AccessToken,
		AccessExpire: tokenResp.AccessExpire,
		RefreshAfter: tokenResp.RefreshAfter,
	}, nil
}
