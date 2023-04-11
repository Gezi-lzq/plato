package source

import (
	"context"

	"github.com/Gezi-lzq/plato/common/config"
	"github.com/Gezi-lzq/plato/common/discovery"
	"github.com/bytedance/gopkg/util/logger"
)

func Init() {
	eventChan = make(chan *Event)
	ctx := context.Background()
	go DataHandler(&ctx)
	if config.IsDebug() {
		ctx := context.Background()
		TestServiceRegiste(&ctx, "7896", "node1")
		TestServiceRegiste(&ctx, "7897", "node2")
		TestServiceRegiste(&ctx, "7898", "node3")
	}
}

func DataHandler(ctx *context.Context) {
	dis := discovery.NewServiceDiscovery(ctx)
	defer dis.Close()
	setFunc := func(key, value string) {
		if ed, err := discovery.UnMarshal([]byte(value)); err == nil {
			if event := NewEvent(ed, AddNodeEvent); ed != nil {
				eventChan <- event
			}
		} else {
			logger.CtxErrorf(*ctx, "DataHandler.setFuc.err:%s", err.Error())
		}
	}
	delFunc := func(key, value string) {
		if ed, err := discovery.UnMarshal([]byte(value)); err == nil {
			if event := NewEvent(ed, DelNodeEvent); ed != nil {
				eventChan <- event
			}
		} else {
			logger.CtxErrorf(*ctx, "DataHandler.delFuc.err:%s", err.Error())
		}
	}
	err := dis.WatchService("/plato/ip_dispatcher", setFunc, delFunc)
	if err != nil {
		panic(err)
	}
}
