package peerjs_exts

import (
	p "github.com/muka/peerjs-go"
	e "github.com/muka/peerjs-go/emitter"
)

type NamedEvents struct {
	Open         func(id string)
	Connection   func(*p.DataConnection)
	Call         func(*p.MediaConnection)
	Close        func(nil any)
	Disconnected func(id string)
	Error        func(p.PeerError)
}

type EventMiddleware func(event string, arg any)

func (nevs *NamedEvents) Join(p *p.Peer, mw EventMiddleware) {
	p.On(wrap(mw, &nevs.Open, "open"))
	p.On(wrap(mw, &nevs.Connection, "connection"))
	p.On(wrap(mw, &nevs.Call, "call"))
	p.On(wrap(mw, &nevs.Close, "close"))
	p.On(wrap(mw, &nevs.Disconnected, "disconnected"))
	p.On(wrap(mw, &nevs.Error, "error"))
}

func wrap[T any](mw EventMiddleware, handler *func(T), event string) (string, e.EventHandler) {
	hd := func(arg any) {
		if h := *handler; h != nil {
			h(arg.(T))
		}
	}
	if mw == nil {
		return event, hd
	} else {
		return event, func(arg any) { mw(event, arg); hd(arg) }
	}
}
