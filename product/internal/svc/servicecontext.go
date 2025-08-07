package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"shop/product/internal/config"
	"shop/product/model"
)

type ServiceContext struct {
	Config       config.Config
	ProductModel model.ProductsModel
	Redis        *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:       c,
		ProductModel: model.NewProductsModel(sqlx.NewMysql(c.Mysql.DataSource)),
		Redis:        redis.MustNewRedis(c.Redis.RedisConf),
	}
}
