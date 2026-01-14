// Package router provides a unified abstraction for routing.
package router

type (
	// Pack should as enbedded meta info for web framework Pack, like EchoPack.
	Pack struct {
		Prefix string // current prefix
	}
	// Route provide abstraction for routing node.
	// P is web framework package transform struct.
	Route[P any] interface {
		// Active report whether current route is active to route and handle.
		Active() bool
		// Path return current route's path, just as ip in network.
		Path() string
		// Handle operations after pack inbound.
		// Explicitly rewriting the Handle call after inbound is necessary when you
		// need the route to do something other than routing.
		Handle(P)
		// Inbound will parse the Pack and update the relevant information.
		Inbound(P)
		// Outbound will send the updated Pack to subsequent nodes.
		Outbound()
		// ToRoute show that next node is a route.
		ToRoute(Route[P])
		// ToPeer show that next node is a perr.
		ToPeer(Peer[P])
	}
	// Peer should collect many endpoints.
	Peer[P any] interface {
		// Parse will parse the Pack context and mount all endpoints.
		Parse(P)
	}
	// Endpoint identify the end point of a processing path.
	Endpoint[F any, MF any] struct {
		Path        string
		Method      string
		Kind        uint // for log or metric
		Handler     F    // handler
		PreHandlers []MF // pre handler
	}
)
