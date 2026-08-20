package svc

import (
	"go-zero-flash-sale/app/usercenter/cmd/rpc/internal/config"
	models "go-zero-flash-sale/app/usercenter/model"

	"go-zero-flash-sale/pkg/token"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config        config.Config
	TxRunner      sqlx.SqlConn
	UserModel     models.UserModel
	UserAuthModel models.UserAuthModel
	Redis         *redis.Redis
	Maker         token.Maker
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DB.DataSource)
	maker, err := token.NewJwtMaker(c.JwtAuth.AccessSecret)
	if err != nil {
		logx.Must(err)
	}
	return &ServiceContext{
		Config:        c,
		TxRunner:      conn,
		UserModel:     models.NewUserModel(conn, c.Cache),
		UserAuthModel: models.NewUserAuthModel(conn, c.Cache),
		Redis:         redis.MustNewRedis(c.Redis.RedisConf),
		Maker:         maker,
	}
}
