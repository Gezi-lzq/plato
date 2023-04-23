package gateway

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"syscall"

	"github.com/Gezi-lzq/plato/common/tcp"

	"github.com/Gezi-lzq/plato/common/config"
)

// RunMain 启动网关服务
func RunMain(path string) {
	config.Init(path)
	ln, err := net.ListenTCP("tcp", &net.TCPAddr{Port: config.GetGatewayServerPort()})
	if err != nil {
		log.Fatalf("StartTCPEPollServer err:%s", err.Error())
		panic(err)
	}
	initWorkPoll()
	initEpoll(ln, runProc)
	fmt.Println("-------------im gateway stated------------")
	select {}
}

func initEpoll(ln *net.TCPListener, f func(c *connection, ep *epoller)) {
	setLimit()
	EpollerPool = newEPool(ln, f)
	EpollerPool.createAcceptProcess()
	EpollerPool.startEPool()
}

func runProc(c *connection, ep *epoller) {
	// step1: 读取一个完整的消息包
	dataBuf, err := tcp.ReadData(c.conn)
	if err != nil {
		// 如果读取conn时发现连接关闭，则直接端口连接
		if errors.Is(err, io.EOF) {
			ep.remove(c)
		}
		return
	}
	err = wPool.Submit(func() {
		// step2:交给 state server rpc 处理
		bytes := tcp.DataPgk{
			Len:  uint32(len(dataBuf)),
			Data: dataBuf,
		}
		tcp.SendData(c.conn, bytes.Marshal())
	})
	if err != nil {
		fmt.Errorf("runProc:err:%+v\n", err.Error())
	}
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
