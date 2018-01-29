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
	"errors"
	"fmt"
	"sync"
)

const defaultNamespace = "/"

type implNamespace struct {
	server       *implServer
	name         string
	connHandlers []func(Socket)
	sockets      *sync.Map // key=sid,value=socket
	rooms        *sync.Map // key=room_name,value=implRoom
}

func (p *implNamespace) ID() string {
	return p.name
}

func (p *implNamespace) OnConnect(callback func(socket Socket)) {
	p.connHandlers = append(p.connHandlers, callback)
}

func (p *implNamespace) GetSockets() Sockets {
	ret := make(Sockets)
	p.sockets.Range(func(key, value interface{}) bool {
		ret[key.(string)] = value.(Socket)
		return true
	})
	return ret
}

func (p *implNamespace) ensureRoom(roomName string) *implRoom {
	var room *implRoom
	if found, ok := p.rooms.Load(roomName); !ok {
		room = newRoom(p, roomName)
		p.rooms.Store(roomName, room)
	} else {
		room = found.(*implRoom)
	}
	return room
}

func (p *implNamespace) appendSocket(socket *implSocket) error {
	p.sockets.Delete(socket.ID())
	return nil
}

func (p *implNamespace) removeSocket(socket *implSocket) (err error) {
	defer func() {
		e := recover()
		if e == nil {
			return
		}
		if p.server.logger.err != nil {
			p.server.logger.err("handle socket create failed: %s\n", e)
		}
		switch e.(type) {
		case error:
			err = e.(error)
			break
		case string:
			err = errors.New(e.(string))
			break
		}
	}()
	sid := socket.ID()
	exist, loaded := p.sockets.LoadOrStore(sid, socket)
	if loaded {
		if exist == socket {
			return nil
		}
		return fmt.Errorf("socket %s exists already", sid)
	}
	for _, fn := range p.connHandlers {
		fn(socket)
	}
	return nil
}

func newNamespace(server *implServer, name string) *implNamespace {
	nsp := implNamespace{
		server:       server,
		name:         name,
		sockets:      new(sync.Map),
		rooms:        new(sync.Map),
		connHandlers: make([]func(Socket), 0),
	}
	return &nsp
}
