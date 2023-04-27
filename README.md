# plato

此项目目前只计划实现接入层部分

![img](doc/img/%E6%9E%B6%E6%9E%84%E8%8D%89%E5%9B%BE.png)

只实现`IP config server`，`gateway`,`state`三个部分

## IP config - 长连接调度引擎

目前已实现 ✔️
![img](/doc/img/ip_conf_server.png)

1. 提供一个查询 endpoint 信息的列表接口 ✔️
2. 设计 评估 机器负载情况 的 评分策略 ✔️
3. 降低在线查询延时时间，提高各节点分值的实时性与有效性 ✔️
4. 可支持水平扩展 ✔️

设计思路：

1. 机器负载的评分策略：

- 分为动态分与静态分两部分

- 由于对于长连接网关，带宽的消耗可以作为一个核心的参考点来作为衡量机器负载的参考。因此可以选取机器的剩余带宽作为动态分的指标。

- 同时长连接网关还存在大量非活跃状态的连接的情况，也会占用较大内存并且需要维护心跳与收发消息协程，因此静态分选用连接数量来衡量负载情况。

- 由于动态分的变化较为活跃，所以需要对数据进行一定处理，需要以 GB 为单位进行近似操作，只保留两位小数。

- 为了降低查询的延迟，提高 ip 列表所体现的负载的排序结果的真实性，故选择通过计算一定时间窗口内的分值，再进行排序。使用 5s 的时间窗口取均值来屏蔽噪声。

- 进行比较过程时，优先比较动态分，其次是静态分。

2. 分值计算：

- ip 列表的查询 与 数据计算与更新 之间 是异步的而非同步关系。

- 仅当请求到达时，会对当前数据进行排序，返回给出 ip 列表

参考资料：

1. [长连接的负载均衡](https://lushunjian.github.io/blog/2018/07/28/%E9%95%BF%E8%BF%9E%E6%8E%A5%E7%9A%84%E8%B4%9F%E8%BD%BD%E5%9D%87%E8%A1%A1/)
2. [长连接负载均衡的问题 | 卡瓦邦噶！](https://www.kawabangga.com/posts/4714)

## Gateway

初步设计：
![img](/doc/img/gateway_v1.png)

具体设计：
![img](/doc/img/gateway.png)

设计目标：

1. 维护长连接 socket 状态 ✔️

2. 协议解析与消息包转发 ✔️

3. 状态统计与上报

设计思路：

1. 关于 C10K 问题，选择使用 epoll 多路复用方案。以此解决静态连接占用应用层内存过多的问题。
2. 使用 Reactor 模式，只有当有读事件或者写事件发生时才会有协程去读写 socket。
3. 多 accept 与 多轮询器 机制 来提高吞吐量。通过全局 channel 来协调 accept 关于 fd 的生产，与 epoller 关于 fd 的消费。
4. 资源池化，对于解析出的 data 交给 work pool 协程池去处理。

关于 state 的拆分设计：

关于 Gateway 的设计的出发点均是为了【提高单机持有连接的数量】与【确保稳定性】

因此在设计的时候需要【尽可能少创建内存结构】

也就是需要对内存结构进行一个“必要”与“仅需要”的等级划分。

- 在 epoll 的网络模型下，socket 内存与维护 socket 的协程的内存均为必要的，不易节省的。

- 创建协程与读取消息的 buffer 对象的内存消耗是必要的，也是可通过池化来约束的。

- 心跳定时器对内存的占用是需要的，但并非必要。

- 维护连接在业务上的键值关系对内存的占用是需要的，但并非必要。

因此，选择将维护 connect 对象与连接相关的状态（clientID，定时器，fid，sessionID...）等等需求，转移到 state server 上。而 gateway 只需进行连接的维护与消息的转发。

参考资料：

[百万 Go TCP 连接的思考: epoll 方式减少资源占用](https://colobu.com/2019/02/23/1m-go-tcp-connection/)

## State

![img](/doc/img/state.png)

1. 使用 unix domin socket 来增大传输速率。

2. 内部维护一个 cmd channel，异步处理 cmd，削峰保证服务稳定。

3. ....

参考资料：

[网络协议之:socket 协议详解之 Unix domain Socket](https://juejin.cn/post/7075509542687080456)

[devlights/go-grpc-uds-example
](https://github.com/devlights/go-grpc-uds-example/blob/master/cmd/server/server.go)
