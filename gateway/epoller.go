package gateway

import (
	"net"
	"reflect"
	"syscall"

	"github.com/Gezi-lzq/plato/common/config"
	"golang.org/x/sys/unix"
)

// epoller 对象 轮询器
type epoller struct {
	fd int
}

func newEpoller() (*epoller, error) {
	fd, err := unix.EpollCreate1(0)
	if err != nil {
		return nil, err
	}
	return &epoller{
		fd: fd,
	}, nil
}

// TODO: 默认水平触发模式,可采用非阻塞FD,优化边沿触发模式
func (e *epoller) add(conn *connection) error {
	// Extract file descriptor associated with the connection
	fd := conn.fd
	err := unix.EpollCtl(e.fd, syscall.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Events: unix.EPOLLIN | unix.EPOLLHUP, Fd: int32(fd)})
	if err != nil {
		return err
	}
	EpollerPool.tables.Store(fd, conn)
	return nil
}

func (e *epoller) remove(conn *connection) error {
	subTcpNum()
	fd := conn.fd
	err := unix.EpollCtl(e.fd, syscall.EPOLL_CTL_DEL, fd, nil)
	if err != nil {
		return err
	}
	EpollerPool.tables.Delete(fd)
	return nil
}

func (e *epoller) wait(msec int) ([]*connection, error) {
	events := make([]unix.EpollEvent, config.GetGatewayEpollWaitQueueSize())
	n, err := unix.EpollWait(e.fd, events, msec)
	if err != nil {
		return nil, err
	}
	var connections []*connection
	for i := 0; i < n; i++ {
		//log.Printf("event:%+v\n", events[i])
		if conn, ok := EpollerPool.tables.Load(int(events[i].Fd)); ok {
			connections = append(connections, conn.(*connection))
		}
	}
	return connections, nil
}

func socketFD(conn *net.TCPConn) int {
	tcpConn := reflect.Indirect(reflect.ValueOf(*conn)).FieldByName("conn")
	fdVal := tcpConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")
	return int(pfdVal.FieldByName("Sysfd").Int())
}

func setTcpConifg(c *net.TCPConn) {
	_ = c.SetKeepAlive(true)
}
