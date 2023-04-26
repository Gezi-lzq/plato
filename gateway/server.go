package gateway

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"syscall"

	"github.com/Gezi-lzq/plato/common/prpc"
	"github.com/Gezi-lzq/plato/common/tcp"
	"github.com/Gezi-lzq/plato/gateway/rpc/client"
	"github.com/Gezi-lzq/plato/gateway/rpc/service"
	"google.golang.org/grpc"

	"github.com/Gezi-lzq/plato/common/config"
)

var cmdChannel chan *service.CmdContext

// RunMain 启动网关服务
func RunMain(path string) {
	config.Init(path)
	ln, err := net.ListenTCP("tcp", &net.TCPAddr{Port: config.GetGatewayTCPServerPort()})
	if err != nil {
		log.Fatalf("StartTCPEPollServer err:%s", err.Error())
		panic(err)
	}
	initWorkPoll()
	initEpoll(ln, runProc)
	fmt.Println("-------------im gateway stated------------")
	// TODO:创建gateway的RPC server 并注册服务
	cmdChannel = make(chan *service.CmdContext, config.GetGatewayCmdChannelNum())
	s := prpc.NewPServer(
		prpc.WithProtocol(config.GetGatewayRPCSockAdd()),
		prpc.WithSockAddr(config.GetGatewayRPCProtocol()),

		prpc.WithServiceName(config.GetGatewayServiceName()),
		prpc.WithIP(config.GetGatewayServiceAddr()),
		prpc.WithPort(config.GetGatewayRPCServerPort()),
		prpc.WithWeight(config.GetGatewayRPCWeight()))
	fmt.Println(config.GetGatewayServiceName(), config.GetGatewayServiceAddr(), config.GetGatewayRPCServerPort(), config.GetGatewayRPCWeight())
	s.RegisterService(func(server *grpc.Server) {
		service.RegisterGatewayServer(server, &service.Service{CmdChannel: cmdChannel})
	})
	// 启动rpc 客户端
	client.Init()
	// 启动 命令处理写协程
	go cmdHandler()
	// 启动 rpc server
	s.Start(context.TODO())
}

func cmdHandler() {
	for cmd := range cmdChannel {
		// 异步提交到协池中完成发送任务
		switch cmd.Cmd {
		case service.DelConnCmd:
			wPool.Submit(func() { closeConn(cmd) })
		case service.PushCmd:
			wPool.Submit(func() { sendMsgByCmd(cmd) })
		default:
			panic("command undefined")
		}
	}
}
func closeConn(cmd *service.CmdContext) {
	if connPtr, ok := EpollerPool.tables.Load(cmd.FD); ok {
		conn, _ := connPtr.(*connection)
		conn.Close()
		EpollerPool.tables.Delete(cmd.FD)
	}
}
func sendMsgByCmd(cmd *service.CmdContext) {
	if connPtr, ok := EpollerPool.tables.Load(cmd.FD); ok {
		conn, _ := connPtr.(*connection)
		dp := tcp.DataPgk{
			Len:  uint32(len(cmd.Playload)),
			Data: cmd.Playload,
		}
		tcp.SendData(conn.conn, dp.Marshal())
	}
}

func initEpoll(ln *net.TCPListener, f func(c *connection, ep *epoller)) {
	setLimit()
	EpollerPool = newEPool(ln, f)
	EpollerPool.createAcceptProcess()
	EpollerPool.startEPool()
}

func runProc(c *connection, ep *epoller) {
	ctx := context.Background() // 起始的contenxt
	// step1: 读取一个完整的消息包
	dataBuf, err := tcp.ReadData(c.conn)
	if err != nil {
		// 如果读取conn时发现连接关闭，则直接端口连接
		// 通知 state 清理掉意外退出的 conn的状态信息
		if errors.Is(err, io.EOF) {
			// 异步操作，不需要等到返回成功在进行
			// 因为消息可靠性的保障是通过协议完成的而非某次cmd
			ep.remove(c)
			client.CancelConn(&ctx, getEndpoint(), int32(c.fd), nil)
		}
		return
	}
	err = wPool.Submit(func() {
		// step2:交给 state server rpc 处理
		client.SendMsg(&ctx, getEndpoint(), int32(c.fd), dataBuf)
	})
	if err != nil {
		fmt.Errorf("runProc:err:%+v\n", err.Error())
	}
}

func getEndpoint() string {
	return fmt.Sprintf("%s:%d", config.GetGatewayServiceAddr(), config.GetGatewayRPCServerPort())
}

// 设置go 进程打开文件数的限制
func setLimit() {
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}
	rLimit.Cur = rLimit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}

	log.Printf("set cur limit: %d", rLimit.Cur)
}
