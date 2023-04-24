package plugin

import (
	"errors"
	"fmt"

	"github.com/Gezi-lzq/plato/common/prpc/config"
	"github.com/Gezi-lzq/plato/common/prpc/discov"
	"github.com/Gezi-lzq/plato/common/prpc/discov/etcd"
)

// GetDiscovInstance 获取服务发现实例
func GetDiscovInstance() (discov.Discovery, error) {
	name := config.GetDiscovName()
	switch name {
	case "etcd":
		return etcd.NewETCDRegister(etcd.WithEndpoints(config.GetDiscovEndpoints()))
	}

	return nil, errors.New(fmt.Sprintf("not exist plugin:%s", name))
}
