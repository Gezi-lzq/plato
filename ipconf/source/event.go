package source

import (
	"fmt"

	"github.com/Gezi-lzq/plato/common/discovery"
)

var eventChan chan *Event

func EventChan() <-chan *Event {
	return eventChan
}

type EvenType string

const (
	AddNodeEvent EvenType = "addNode"
	DelNodeEvent EvenType = "delNode"
)

type Event struct {
	Type         EvenType
	IP           string
	Port         string
	ConnectNum   float64
	MessageBytes float64
}

func NewEvent(ed *discovery.EndpointInfo, eventType EvenType) *Event {
	if ed == nil || ed.MetaData == nil {
		return nil
	}
	var connNum, msgBytes float64
	if data, ok := ed.MetaData["connect_num"]; ok {
		connNum = data.(float64) // 若出错，此处应panic 暴露错误
	}
	if data, ok := ed.MetaData["message_bytes"]; ok {
		msgBytes = data.(float64) // 若出错，此处应panic 暴露错误
	}
	return &Event{
		Type:         eventType,
		IP:           ed.IP,
		Port:         ed.Port,
		ConnectNum:   connNum,
		MessageBytes: msgBytes,
	}
}

func (nd *Event) Key() string {
	return fmt.Sprintf("%s:%s", nd.IP, nd.Port)
}
