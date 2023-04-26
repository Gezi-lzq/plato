package discov

import (
	"fmt"
	"strconv"
	"strings"
)

type Service struct {
	Name      string      `json:"name"`
	Endpoints []*Endpoint `json:"endpoints"`
}

type Endpoint struct {
	ServerName string `json:"server_name"`
	Protocol   string `json:protocal`
	SockAddr   string `json:sockAddr`
	IP         string `json:"ip"`
	Port       int    `json:"port"`
	Weight     int    `json:"weight"`
	Enable     bool   `json:"enable"`
}

func (endpoint *Endpoint) GetAddr() string {
	switch endpoint.Protocol {
	case "tcp", "tcp4", "tcp6":
		return fmt.Sprintf("%s:%d", endpoint.IP, endpoint.Port)
	case "unix", "unixgram", "unixpacket":
		return fmt.Sprintf("unix:///%s", endpoint.SockAddr)
	default:
		return ""
	}
}

func InitEndpointByAddr(addr string) *Endpoint {
	addrs := strings.Split(addr, ":")
	if len(addrs) == 2 {
		port, _ := strconv.Atoi(addrs[1])
		return &Endpoint{
			Protocol: "tcp",
			IP:       addrs[0],
			Port:     port,
		}
	} else {
		return &Endpoint{
			Protocol: "unix",
			SockAddr: addr,
		}
	}
}
