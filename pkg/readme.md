---
title: "pkg section"
date: "2026/01/05 pm"
author: "wendisx"
---

# pkg section

- [database](#db)
- [log](#log)
- [router](#router)
- [config](#config)

---

some simple design ideas.

```go
/*
    索性直接用golang code block来写思路. 实际上, puzzle我有一段时间没有编写了. 最早是为了自己能够在闲暇之余能够找点除了video之外的乐趣, 比如写点通用的代码, 做点通用的工程. 作为一个主要做backend的人, 当然需要写好一个差不多的server框架才行了. (这里并不是直接写web framework, 那个没时间研究同时没技术完成.) 这次不写那么多废话了, 毕竟是写给自己使用的, 怎么方便怎么来. 也是考完试了, 写点golang放松一下, 总不能每天写算法. (虽然算法水平一般.)
*/
// 防止因为过多的随意取名导致function调用时来回切换, 简单介绍一些原则.
func InitX(args...) // 初始化默认的对象
func SetupX(args...) // 配置默认对象: 采用function配置形式
func ResetX() // 重置对象属性
func FreeX() // 释放对象
```

## <a id="db">database</a>

database 如何设计api. 其实这就是调用习惯的问题: init() -> setup() -> run() -> free(). database在后端项目中的唯一作用就是维护一个local database buffer pool用于直接和database进行交互. 为了软件工程上的开闭原则, 最小化的对象操作在database层面上无非就是crud四个. 因此设计时只需要尝试针对一个类型的database设计出对指定对象的最小可抽象操作即可.

```go
/*
    database -- 统一数据库层的抽象, 简化集成和操作
    loc(package): pkg/db
       loc(test): pkg/db/database_test [passed]
*/
package database
```

## <a id="log">log</a>

log 设计主要包含几个要点: 显示提示, 内容输出, 结构化和非结构化. 在 container 环境中, 显然最核心的就是非结构化日志, 主要用于人可读的形式呈现, 而同时又需要保留结构化接口, 用于日志分析. 因此直接封装 slog 可以合理地处理这些. 

```go
/*
    clog -- colorable log, 封装自slog设计
    loc(package): pkg/clog
       loc(test): pkg/clog/clog_test [passed]
*/
package clog
```

## <a id="router">router</a>

router 是一个极其简单但又极其复杂的设计部分, 不过在一定程度上其实可以直接借鉴许多网络概念设计, 例如l3层的路由器处理package的行为.

详细记录一下router的设计, 因为这在一定程度上确实有点复杂. 但是对于熟悉网络的人而言, 这其实很简单. 我们都知道router其实更像一棵树, 从/出发走到最终的端点路径, 这最终会组合为一条简单的uri就像`/agent/api/v1/namespace/endpoints...`, 这只是一个举例, 但是一般的uri会以某种类似的格式划分, 但是其本质上都是: 从/出发, 真的如此吗, 文件系统是如此的, 但是你得知道1必须从0开始, 我们将所有的路径按照`/`划分, 最终得到的就是一些以`/`开头的path, 这将成为route的一种标记, 我们希望这可以在一定的局部唯一, 因为可能可以见到`/api/api`这样的奇怪路由. 接下来看一些设定:

P: Pack - 包, 代表网络上的传递上下文.</br>
Route[P]:  传递指定Pack的路由, 用于明确职责. 代表网络中可以接收和转发操作的中继器.</br>
Peer[P]: 端点, 代表一段链路的终点.</br>
Endpoint[F,MF]: 端点处理器, 指明关联的Peer[P]允许做的事情.

```go
/*
    router -- 网络路由抽象设计, 简化统一路由层抽象
    loc(package): pkg/router 
       loc(test): pkg/router/router_test [passed]
*/
package router
type (
  Pack struct {
    Prefix string // 记录前缀路径 
    // 其他framework无关性上下文
  }
  Route[P any] interface {
    Active() bool // 是否工作
    Path() string // 终止递归
    Inbound(P) // 路由接收
    Outbound() // 路由转发
    ToRoute(Route[P]) // next route
    ToPeer(Peer[P]) // next peer
    Handle(P) // route 执行操作
  }
  Peer[P any] interface {
    Parse(P) // 解析
  }
  Endpoint[F any, MF any] struct {
    Method string
    Path string
    Handler F
    PreHandlers []MF
  }
)

/* 
    1. 封装 Pack 带分组器, 存在构造 Pack 的public api, 这里只能允许自行构建原始 Pack.
    2. 封装 Route[P] 保存 Pack上下文, 自身路径, next route和next peer. 需要手动构造, 对封装后的实际 Route进行接口实现
      2.1 func (r *_(Route[P])) Active() bool {return true} <- 这里一定返回 true, 默认启动
      2.2 func (~) Inbound(p Pack) { ... <- 这里一定需要实现的操作为更新Pack上下文记录当前路径到prefix后进行grouper切分} // 由于默认Handle()为警告提醒显式实现, 该方法在需要执行Handle()时需要显式重写.
      2.3 func (~) Outbound() { ... <- 这里一定需要遍历所有的next route和next peer进行上下文传递; 为了使外部重写的Active()覆盖内部的Active(), 不能等待到内部的Route转换后在Inbound()中进行拦截(实际上Inbound()也判断了), 在遍历next route时, 需要确保下一个有效后才传递上下文}
      2.4 递归调用Outbound()进行全链路执行
    3. 封装 Peer[P] 保存next endpoint. 无需构造, 直接定义, 实现Parse(P) 和 ToEndpoint(...) 即可. 内部的列表可以在触发ToEndpoint()时构造. 重写最外部的Parse(P), 定义所有端点后调用内部的Parse(P)传递上下文即可,
    4. 构造多个嵌入封装, 组合, 最终执行gateway的Outbound()即可开始路由.
*/

```

## <a id="config">config</a>

config 应该是简单的一个环节, 主要的功能就是加载必要的配置选项后进行顺序化处理得到最终的配置并装载到组件当中, 可以看作是组件行为的一般约束接口(但是受限). 无法尝试对已经实现的第三方组件进行更多层面的配置, 但是可以想到的是通过参数改变组件的行为表现. 一般性config的加载顺序大多为: `config file -> environment variable -> cli flags`, 采用这种加载顺序的一大最重要的原因就是: 后面的配置不需要判定可以直接覆盖前面的配置, 同时在不存在时可以默认使用前面的配置. config file加载时提供了一个最为全面的模板用于记录所有可能的参数, environment variable只对存在的映射进行修改, 最终的flags也是通过映射进行修改.  

```go
/*
    config -- 加载和处理配置, 提供便于操作配置的接口
    loc(package): pkg/config
       loc(test): 
         - pkg/config/datadict_test [passed]
         - pkg/config/config_test [passed]
*/
package config
```

