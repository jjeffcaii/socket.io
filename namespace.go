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
	"sync"
)

const defaultNamespace = "/"

type implNamespace struct {
	server       *implServer
	name         string
	connHandlers []func(Socket)
	sockets      map[string]*implSocket
	locker       *sync.RWMutex
}

func (p *implNamespace) ID() string {
	return p.name
}

func (p *implNamespace) OnConnect(callback func(socket Socket)) {
	p.connHandlers = append(p.connHandlers, callback)
}

func (p *implNamespace) GetSockets() Sockets {
	ret := make(Sockets)
	p.locker.RLock()
	for k, v := range p.sockets {
		ret[k] = v
	}
	p.locker.RUnlock()
	return ret
}

func (p *implNamespace) leaveSocket(socket *implSocket) error {
	p.locker.Lock()
	defer p.locker.Unlock()
	delete(p.sockets, socket.ID())
	return nil
}

func (p *implNamespace) joinSocket(socket *implSocket) error {
	sid := socket.ID()
	p.locker.Lock()
	exist, ok := p.sockets[sid]
	p.locker.Unlock()
	if !ok {
		p.sockets[sid] = socket
	} else if exist == socket {
		return nil
	} else {
		return fmt.Errorf("socket %s exists already", sid)
	}
	c := len(p.connHandlers)
	if c < 1 {
		return nil
	}

	if c == 1 {
		p.connHandlers[0](socket)
		return nil
	}

	wg := new(sync.WaitGroup)
	wg.Add(len(p.connHandlers))
	for _, fn := range p.connHandlers {
		go func() {
			defer func() {
				if e := recover(); e != nil {
					if p.server.logger.err != nil {
						p.server.logger.err.Println("handle socket create failed:", e)
					}
				}
				wg.Done()
			}()
			fn(socket)
		}()
	}
	wg.Wait()
	return nil
}

func newNamespace(server *implServer, name string) *implNamespace {
	nsp := implNamespace{
		server:       server,
		name:         name,
		sockets:      make(map[string]*implSocket),
		connHandlers: make([]func(Socket), 0),
		locker:       new(sync.RWMutex),
	}
	return &nsp
}
