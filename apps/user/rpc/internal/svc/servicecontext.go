package svc

import (
	"SAI-IM/apps/user/models"
	"SAI-IM/apps/user/rpc/internal/config"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	models.UsersModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.Mysql.DataSource)

	return &ServiceContext{
		Config: c,

		UsersModel: models.NewUsersModel(sqlConn, c.Cache),
	}
}
