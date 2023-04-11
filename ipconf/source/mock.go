package source

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/Gezi-lzq/plato/common/config"
	"github.com/Gezi-lzq/plato/common/discovery"
)

func TestServiceRegiste(ctx *context.Context, port, node string) {
	// 模拟服务发现
	go func() {
		ed := discovery.EndpointInfo{
			IP:   "127.0.0,1",
			Port: port,
			MetaData: map[string]interface{}{
				"connect_num":   float64(rand.Int63n(12312312312123)),
				"message_bytes": float64(rand.Int63n(123123123123)),
			},
		}
		sr, err := discovery.NewServiceRegister(ctx,
			fmt.Sprintf("%s/%s", config.GetServicePathForIPConf(), node), &ed, time.Now().Unix())
		if err != nil {
			panic(err)
		}
		go sr.ListenLeaseRespChan()
		for {
			ed = discovery.EndpointInfo{
				IP:   "17.0.0.1",
				Port: port,
				MetaData: map[string]interface{}{
					"connect_num":   float64(rand.Int63n(12312312312123123)),
					"message_bytes": float64(rand.Int63n(12312312313116)),
				},
			}
			sr.UpdateValue(&ed)
			time.Sleep(1 * time.Second)
		}
	}()
}
