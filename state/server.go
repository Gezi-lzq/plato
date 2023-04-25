package state

import (
	"context"
	"fmt"

	"github.com/Gezi-lzq/plato/common/config"
	"github.com/Gezi-lzq/plato/common/prpc"
	"github.com/Gezi-lzq/plato/state/rpc/service"
	"google.golang.org/grpc"
)

var cmdChannel chan *service.CmdContext

// RunMain 启动state服务
func RunMain(path string) {
	config.Init(path)
	cmdChannel = make(chan *service.CmdContext, config.GetStateCmdChannelNum())

	s := prpc.NewPServer(
		prpc.WithServiceName(config.GetStateServiceName()),
		prpc.WithIP(config.GetStateServiceAddr()),
		prpc.WithPort(config.GetStateServerPort()),
		prpc.WithWeight(config.GetStateRPCWeight()))
	s.RegisterService(func(server *grpc.Server) {
		service.RegisterStateServer(server, &service.Service{CmdChannel: cmdChannel})
	})
	// 初始化RPC 客户端

	// 启动 命令处理写协程
	go cmdHandler()
	// 启动 rpc server
	s.Start(context.TODO())
}

func cmdHandler() {
	for cmd := range cmdChannel {
		switch cmd.Cmd {
		case service.CancelConnCmd:
			fmt.Printf("cancelconn endpoint:%s, fd:%d, data:%+v", cmd.Endpoint, cmd.FD, cmd.Playload)
		case service.SendMsgCmd:
			fmt.Println("cmdHandler", string(cmd.Playload))
		}
	}
}
