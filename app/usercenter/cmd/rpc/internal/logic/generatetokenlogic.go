package logic

import (
	"context"
	"time"

	"go-zero-flash-sale/app/usercenter/cmd/rpc/internal/svc"
	"go-zero-flash-sale/app/usercenter/cmd/rpc/usercenter"
	models "go-zero-flash-sale/app/usercenter/model"
	"go-zero-flash-sale/pkg/apperror"
	"go-zero-flash-sale/pkg/constants"
	"go-zero-flash-sale/pkg/token"

	"github.com/zeromicro/go-zero/core/logx"
)

type GenerateTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGenerateTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateTokenLogic {
	return &GenerateTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GenerateTokenLogic) GenerateToken(in *usercenter.GenerateTokenRequest) (*usercenter.GenerateTokenResponse, error) {
	// 获取 user
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, in.UserId)
	if err != nil {
		if err == models.ErrNotFound {
			return nil, apperror.ErrUserNotFound
		}
		return nil, apperror.WrapSrvErrf(err, "铸造 token 时出错, user: %v", user)
	}

	// 铸造 token
	accessToken, expiredAt, refreshAfter, err := generateToken(user.Id, l.svcCtx.Config.JwtAuth.AccessExpire, l.svcCtx.Maker)
	if err != nil {
		return nil, err
	}

	return &usercenter.GenerateTokenResponse{
		AccessToken:  accessToken,
		AccessExpire: expiredAt,
		RefreshAfter: refreshAfter,
	}, nil
}

func generateToken(userId, accessExpire int64, maker token.Maker) (string, int64, int64, error) {
	// 铸造 token
	accessToken, payload, err := maker.CreateToken(userId,
		time.Duration(accessExpire),
		token.TokenTypeAccessToken)
	if err != nil {
		return "", constants.InvalidTimeStamp, constants.InvalidTimeStamp, apperror.WrapSrvErrf(err, "CreateToken 调用失败, userId: %d", userId)
	}

	// 计算提醒时间
	issuedAt, expiredAt := payload.IssuedAt.Unix(), payload.ExpiredAt.Unix()
	refreshAfter := issuedAt + (expiredAt-issuedAt)/2

	return accessToken, expiredAt, refreshAfter, nil
}
