package domain

import (
	"sync/atomic"
	"unsafe"
)

type Endpoint struct {
	IP          string       `string:"ip"`
	Port        string       `json:"port"`
	ActiveSorce float64      `json:"-"`
	StaticSorce float64      `json:"-"`
	Stats       *Stat        `json:"-"`
	window      *stateWindow `json:"-"`
}

func NewEndpoint(ip, port string) *Endpoint {
	ed := &Endpoint{
		IP:   ip,
		Port: port,
	}
	ed.window = newStateWindow()
	ed.Stats = ed.window.getStat()
	go func() {
		for stat := range ed.window.statChan {
			ed.window.appendStat(stat)
			newStat := ed.window.getStat()
			// 将newStat更新到 ed.Stats 指针所指向的内存地址中(指针类型的原子性赋值操作)
			// 1. 表示将 ed.Stats 变量的地址转换为 unsafe.Pointer 类型的指针，
			//    并将该指针再次转换为一个指针类型的指针（即指向 unsafe.Pointer 类型的指针的指针）
			// 2. 表示将 newStat 变量的地址转换为 unsafe.Pointer 类型的指针。
			atomic.SwapPointer((*unsafe.Pointer)((unsafe.Pointer)(ed.Stats)), unsafe.Pointer(newStat))
		}
	}()
	return ed
}

func (ed *Endpoint) UpdateStat(s *Stat) {
	ed.window.statChan <- s
}

func (ed *Endpoint) CalculateScore(stx *IpConfContext) {
	// 如何 stats字段为空，则直接使用上一次计算结果，此次不更新
	if ed.Stats != nil {
		ed.ActiveSorce = ed.Stats.CalculateActiveSorce()
		ed.StaticSorce = ed.Stats.CalculateStaticSorce()
	}
}
