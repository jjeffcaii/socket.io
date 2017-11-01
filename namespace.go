package sio

import "sync"

type implNamespace struct {
	server       Server
	name         string
	sockets      *sync.Map
	connHandlers []func(Socket)
}

func (p *implNamespace) OnConnect(callback func(socket Socket)) {
	panic("implement me")
}

func (p *implNamespace) GetSockets() Sockets {
	ret := make(Sockets)
	p.sockets.Range(func(key, value interface{}) bool {
		k := key.(string)
		v := value.(Socket)
		ret[k] = v
		return false
	})
	return ret
}

func newNamespace(server Server, name string) Namespace {
	nsp := implNamespace{
		server:  server,
		name:    name,
		sockets: new(sync.Map),
	}
	return &nsp
}
