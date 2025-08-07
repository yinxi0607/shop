package svc

import (
	//"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"shop/user/internal/config"
	"shop/user/model"
)

type ServiceContext struct {
	Config    config.Config
	UserModel model.UsersModel
	Redis     *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		UserModel: model.NewUsersModel(sqlx.NewMysql(c.Mysql.DataSource)),
		Redis:     redis.MustNewRedis(c.Redis.RedisConf),
	}
}
