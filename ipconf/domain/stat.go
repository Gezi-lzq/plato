package domain

import "math"

// 对于gateway网关来说，存在不同时期加入进来的物理机，所以机器的配置是不同的，使用负载来衡量会导致偏差。
// 为更好的应对动态的机器配置变化，我们统计其剩余资源值，来衡量一个机器其是否更适合增加其负载
// 这里的数值代表的是，此endpoint对应的机器，其自身剩余的资源指标
type Stat struct {
	ConnectNum   float64
	MessageBytes float64
}

// 因此假设网络带宽是系统瓶颈所在，那么哪台机器富裕的带宽资源多，哪台机器的负载就是最轻的。
// TODO：如何评估他的数量级？何时使用静态值衡量
// TODO：json的解析失败
func (s *Stat) CalculateActiveSorce() float64 {
	return getGB(s.MessageBytes)
}

func (s *Stat) CalculateStaticSorce() float64 {
	return s.ConnectNum
}

// 将以字节为单位转换成以千兆字节（GB）为单位
func getGB(m float64) float64 {
	return decimal(m / (1 << 30))
}

// 将给定的浮点数 value 转换为保留两位小数的浮点数(四舍五入)
func decimal(value float64) float64 {
	return math.Trunc(value*1e2+0.5) * 1e-2
}

func (s *Stat) Clone() *Stat {
	newStat := &Stat{
		MessageBytes: s.MessageBytes,
		ConnectNum:   s.ConnectNum,
	}
	return newStat
}

func (s *Stat) Avg(num float64) {
	s.ConnectNum /= num
	s.MessageBytes /= num
}

func (s *Stat) Sub(st *Stat) {
	if st == nil {
		return
	}
	s.ConnectNum -= st.ConnectNum
	s.MessageBytes -= st.MessageBytes
}

func (s *Stat) Add(st *Stat) {
	if st == nil {
		return
	}
	s.ConnectNum += st.ConnectNum
	s.MessageBytes += st.MessageBytes
}
