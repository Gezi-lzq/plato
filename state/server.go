package state

import (
	"context"
	"fmt"

	"github.com/Gezi-lzq/plato/common/config"
	"github.com/Gezi-lzq/plato/common/prpc"
	"github.com/Gezi-lzq/plato/state/rpc/client"
	"github.com/Gezi-lzq/plato/state/rpc/service"
	"google.golang.org/grpc"
)

var cmdChannel chan *service.CmdContext

// RunMain 启动state服务
func RunMain(path string) {
	config.Init(path)
	cmdChannel = make(chan *service.CmdContext, config.GetStateCmdChannelNum())

	s := prpc.NewPServer(
		prpc.WithSockAddr(config.GetStateRPCProtocol()),
		prpc.WithProtocol(config.GetStateRPCSockAdd()),

		prpc.WithServiceName(config.GetStateServiceName()),
		prpc.WithWeight(config.GetStateRPCWeight()))
	s.RegisterService(func(server *grpc.Server) {
		service.RegisterStateServer(server, &service.Service{CmdChannel: cmdChannel})
	})
	// 初始化RPC 客户端
	client.Init()
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
			fmt.Println("cmdHandler", int32(cmd.FD), string(cmd.Playload))
			client.Push(cmd.Ctx, int32(cmd.FD), cmd.Playload)
		}
	}
}
