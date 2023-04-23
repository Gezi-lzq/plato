package sdk

import (
	"encoding/json"
	"net"

	"github.com/Gezi-lzq/plato/common/tcp"
)

type connect struct {
	sendChan, recvChan chan *Message
	conn               *net.TCPConn
}

// 初始化过程中，启动一个监听协程，将消息解析传入recvChan中
func newConnet(ip net.IP, port int) *connect {
	clientConn := &connect{
		sendChan: make(chan *Message),
		recvChan: make(chan *Message),
	}
	addr := &net.TCPAddr{IP: ip, Port: port}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		panic(err)
	}
	clientConn.conn = conn
	go func() {
		for {
			data, err := tcp.ReadData(conn)
			if err != nil {
				panic(err)
			}
			msg := &Message{}
			err = json.Unmarshal(data, msg)
			if err != nil {
				panic(err)
			}
			clientConn.recvChan <- msg
		}
	}()
	return clientConn
}

func (c *connect) send(data *Message) {
	// 直接发送给接收方
	c.recvChan <- data
}

func (c *connect) recv() <-chan *Message {
	return c.recvChan
}

func (c *connect) close() {
	// 目前没啥值得回收的
}
