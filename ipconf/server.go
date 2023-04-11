package ipconf

import (
	"github.com/Gezi-lzq/plato/common/config"
	"github.com/Gezi-lzq/plato/ipconf/domain"
	"github.com/Gezi-lzq/plato/ipconf/source"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func RunMain(path string) {
	// 初始化加载配置文件
	config.Init(path)
	// 数据源优先启动
	source.Init()
	// 初始化调度层
	domain.Init()
	s := server.Default(server.WithHostPorts(":6789"))
	s.GET("/ip/list", GetIpInfoList)
	s.Spin()
}
