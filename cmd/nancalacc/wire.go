//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"nancalacc/internal/biz"
	"nancalacc/internal/conf"
	"nancalacc/internal/data"
	"nancalacc/internal/server"
	"nancalacc/internal/service"
	"nancalacc/internal/task"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(
		server.ProviderSet,  // HTTP/gRPC服务
		data.ProviderSet,    // 数据层（含DB、第三方API）
		biz.ProviderSet,     // 业务逻辑层
		service.ProviderSet, // 服务层
		task.ProviderSet,    // 后台任务
		newApp))             // 应用入口
}
