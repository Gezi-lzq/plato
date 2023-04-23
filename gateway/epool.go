package gateway

import (
	"fmt"
	"net"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"

	"github.com/Gezi-lzq/plato/common/config"
)

var (
	EpollerPool *ePool // epoll池
	tcpNum      int32  // 当前服务允许接入的最大tcp连接数
)

type ePool struct {
	eChan  chan *connection
	tables sync.Map
	eSize  int
	done   chan struct{}

	ln *net.TCPListener
	f  func(c *connection, ep *epoller)
}

func newEPool(ln *net.TCPListener, cb func(c *connection, ep *epoller)) *ePool {
	return &ePool{
		eChan:  make(chan *connection),
		done:   make(chan struct{}),
		eSize:  config.GetGatewayEpollerNum(),
		tables: sync.Map{},
		ln:     ln,
		f:      cb,
	}
}

func (e *ePool) createAcceptProcess() {
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				conn, e := e.ln.AcceptTCP()
				if e != nil {
					if ne, ok := e.(net.Error); ok && ne.Temporary() {
						fmt.Errorf("accept temp err:%v \n", ne)
						continue
					}
					fmt.Errorf("accept err: %v \n", e)
				}
				// 限流熔断
				if !checkTcp() {
					_ = conn.Close()
					continue
				}
				setTcpConifg(conn)
				c := connection{
					conn: conn,
					fd:   socketFD(conn),
				}
				EpollerPool.addTask(&c)
			}
		}()
	}
}

func (e *ePool) addTask(c *connection) {
	e.eChan <- c
}

func (e *ePool) startEPool() {
	for i := 0; i < e.eSize; i++ {
		go e.startEProc()
	}
}

// 轮询器池 处理器
func (e *ePool) startEProc() {
	ep, err := newEpoller()
	if err != nil {
		panic(err)
	}
	// 监听连接创建事件
	go func() {
		for {
			select {
			case <-e.done:
				return
			case conn := <-e.eChan:
				if err := ep.add(conn); err != nil {
					fmt.Printf("failed to add connection %v \n", err)
					conn.Close() // 直接关闭连接
					continue
				}
				addTcpNum()
				fmt.Printf("EpollerPool new connection[%v] tcpSize:%d \n", conn.conn.RemoteAddr(), tcpNum)
			}
		}
	}()
	// 轮询器在这里轮询等待, 当有wait发生时则调用回调函数去处理
	for {
		select {
		case <-e.done:
			return
		default:
			connections, err := ep.wait(200) // 200ms 一次轮询避免 忙轮询

			if err != nil && err != syscall.EINTR {
				fmt.Printf("failed to epoll wait %v\n", err)
				continue
			}
			for _, conn := range connections {
				if conn == nil {
					break
				}
				e.f(conn, ep)
			}
		}
	}
}

func addTcpNum() {
	atomic.AddInt32(&tcpNum, 1)
}

func getTcpNum() int32 {
	return atomic.LoadInt32(&tcpNum)
}
func subTcpNum() {
	atomic.AddInt32(&tcpNum, -1)
}

func checkTcp() bool {
	num := getTcpNum()
	maxTcpNum := config.GetGatewayMaxTcpNum()
	return num <= maxTcpNum
}
