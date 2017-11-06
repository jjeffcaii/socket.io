// MIT License
//
// Copyright (c) 2017 jjeffcaii@outlook.com
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
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

func newServer(engine eio.Engine) *implServer {
	serv := &implServer{
		engine:     engine,
		namespaces: make(map[string]*implNamespace),
		locker:     new(sync.RWMutex),
	}
	engine.OnConnect(func(rawSocket eio.Socket) {
		newSocket(serv, rawSocket)
	})
	return serv
}
