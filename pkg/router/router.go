package router

/*
	router
	不同 web framework 实现的方式差异很大, 因此不能简单地直接将一个uniform的东西设计出来接管所有的框架
	但是可以设计一套合理的接口将不同框架的差异减小
	调用过程:
	root_pack -> _(Route)_N -> _(Peer)_N -> Mount() -> endpoints
	1. 创建 根Pack定义, 作为路径上下文传递
	2. 封装 Route[Pack], 通过Inbound()接入根Pack, 通过ToRoute()接入next route
	3. 封装 Peer[Pack], 通过ToPeer()接入指定Route下
	4. 重写需要覆盖的中间件在Handle()中, 通过重写Inbound()最后执行Handle()
	5. 重写默认的Mount()加入Endpoints, 最终执行默认Mount()
	6. 执行带根Inbound的Outbound()
*/

type (
	// Pack should as enbedded meta info
	Pack struct {
		Prefix string
		// ...some message to transform
	}
	// P is web framework package transform struct
	Route[P any] interface {
		Active() bool // just as the thing to block some bad or discard route link.
		Path() string // get current path
		Handle(P)     // just add all handlers here
		Inbound(P)
		Outbound()
		ToRoute(Route[P])
		ToPeer(Peer[P])
	}
	// maybe Peer should collect many endpoints.
	Peer[P any] interface {
		Parse(P)
	}
	Endpoint[F any, MF any] struct {
		Path        string
		Method      string
		Kind        uint
		Handler     F
		PreHandlers []MF
	}
)

func init() {
}
