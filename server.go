package sio

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/jjeffcaii/engine.io"
)

type serverOptions struct {
}

type implServer struct {
	engine     eio.Engine
	namespaces map[string]Namespace
	locker     *sync.RWMutex
}

func (p *implServer) Of(nsp string) (Namespace, error) {
	if !isValidNamespace(nsp) {
		return nil, fmt.Errorf("invalid namespace: %s", nsp)
	}
	p.locker.Lock()
	defer p.locker.Unlock()
	if _, ok := p.namespaces[nsp]; ok {
		return nil, fmt.Errorf("namespace %s exists already", nsp)
	}
	n := newNamespace(p, nsp)
	p.namespaces[nsp] = n
	return n, nil
}

func (p *implServer) Router() func(http.ResponseWriter, *http.Request) {
	return p.engine.Router()
}

func (p *implServer) GetSockets() Sockets {
	panic("implement me")
}

func (p *implServer) Close() {
	panic("implement me")
}

func newServer(engine eio.Engine) *implServer {
	s := implServer{
		engine:     engine,
		namespaces: make(map[string]Namespace),
		locker:     new(sync.RWMutex),
	}
	engine.OnConnect(func(socket eio.Socket) {
		socket.OnMessage(func(data []byte) {



		})
	})
	return &s
}
