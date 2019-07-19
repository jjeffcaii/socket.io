package sio

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/jjeffcaii/engine.io"
)

type loggerGroup struct {
	info, warn, err func(format string, v ...interface{})
}

type implServer struct {
	logger     *loggerGroup
	engine     eio.Engine
	namespaces map[string]*implNamespace
	locker     *sync.RWMutex
}

func (p *implServer) Of(nsp string) Namespace {
	if !isValidNamespace(nsp) {
		panic(fmt.Errorf("invalid namespace: %s", nsp))
	}
	p.locker.Lock()
	defer p.locker.Unlock()
	if exist, ok := p.namespaces[nsp]; ok {
		return exist
	}
	foo := newNamespace(p, nsp)
	p.namespaces[nsp] = foo
	return foo
}

func (p *implServer) Router() func(http.ResponseWriter, *http.Request) {
	return p.engine.Router()
}

func (p *implServer) GetSockets() Sockets {
	ret := make(map[string]Socket)
	nsps := make([]*implNamespace, 0)
	p.locker.RLock()
	for _, nsp := range p.namespaces {
		nsps = append(nsps, nsp)
	}
	p.locker.RUnlock()
	for _, nsp := range nsps {
		for k, v := range nsp.GetSockets() {
			ret[k] = v
		}
	}
	return ret
}

func (p *implServer) loadNamespace(nsp string) (*implNamespace, bool) {
	p.locker.RLock()
	defer p.locker.RUnlock()
	if nsp == "" {
		n, ok := p.namespaces[defaultNamespace]
		return n, ok
	}
	n, ok := p.namespaces[nsp]
	return n, ok
}

func newServer(engine eio.Engine, l1, l2, l3 func(format string, v ...interface{})) (s *implServer) {
	s = &implServer{
		engine:     engine,
		namespaces: make(map[string]*implNamespace),
		locker:     new(sync.RWMutex),
		logger: &loggerGroup{
			info: l1,
			warn: l2,
			err:  l3,
		},
	}
	engine.OnConnect(func(rawSocket eio.Socket) {
		// TODO: temporary fix
		go handleSocket(s, rawSocket)
	})
	return
}
