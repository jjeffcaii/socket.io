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

	"github.com/jjeffcaii/engine.io"
	"github.com/jjeffcaii/socket.io/parser"
)

var initData = "0"

type eventHandlers map[string][]func(Message)

func (s eventHandlers) register(event string, handler func(Message)) error {
	if exists, ok := s[event]; ok {
		s[event] = append(exists, handler)
	} else {
		s[event] = []func(Message){handler}
	}
	return nil
}

func (s eventHandlers) exec(event string, msg Message) {
	handlers, ok := s[event]
	if !ok {
		return
	}
	for _, fn := range handlers {
		fn(msg)
	}
}

type implSocket struct {
	nsp           *implNamespace
	conn          eio.Socket
	eventHandlers eventHandlers
	handshake     *Handshake
}

func (p *implSocket) To(room string) ToRoom {
	panic("implement me")
}

func (p *implSocket) In(room string) InRoom {
	panic("implement me")
}

func (p *implSocket) Join(room string) Socket {
	p.nsp.ensureRoom(room).join(p)
	return p
}

func (p *implSocket) Leave(room string) Socket {
	p.nsp.ensureRoom(room).leave(p)
	return p
}

func (p *implSocket) LeaveAll() error {
	panic("implement me")
}

func (p *implSocket) Handshake() *Handshake {
	return p.handshake
}

func (p *implSocket) Namespace() Namespace {
	return p.nsp
}

func (p *implSocket) Emit(event string, first interface{}, others ...interface{}) error {
	if p.nsp == nil {
		return fmt.Errorf("socket is closing")
	}
	body := []interface{}{first}
	if len(others) > 0 {
		body = append(body, others...)
	}
	evt := &parser.MEvent{
		Namespace: p.nsp.name,
		Event:     event,
		Data:      body,
	}
	packet, err := evt.ToPacket()
	if err != nil {
		return err
	}
	bs, err := parser.Encode(packet)
	if err != nil {
		return err
	}
	switch packet.Type {
	default:
		return fmt.Errorf("unsupport emit packet type: %d", packet.Type)
	case parser.EVENT:
		return p.conn.Send(string(bs))
	case parser.EVENTBIN:
		return p.conn.Send(bs)
	}
}

func (p *implSocket) On(event string, callback func(msg Message)) error {
	return p.eventHandlers.register(event, callback)
}

func (p *implSocket) OnError(callback func(error)) error {
	panic("implement me")
}

func (p *implSocket) OnClose(callback func(reason string)) error {
	p.conn.OnClose(callback)
	return nil
}

func (p *implSocket) Close() {
	p.conn.Close()
}

func (p *implSocket) ID() string {
	return p.conn.ID()
}

func (p *implSocket) accept(evt *parser.MEvent) {
	handlers, ok := p.eventHandlers[evt.Event]
	if !ok {
		if p.nsp.server.logger.warn != nil {
			p.nsp.server.logger.warn("no such event '%s'\n", evt.Event)
		}
		return
	}
	for _, it := range evt.Data {
		for _, fn := range handlers {
			fn(it)
		}
	}
}

func newHandshake(req *http.Request) *Handshake {
	return &Handshake{
		Headers: req.Header,
		Query:   req.URL.Query(),
		URL:     req.URL.Path,
		Address: req.RemoteAddr,
	}
}

func newSocket(server *implServer, conn eio.Socket) *implSocket {
	socket := &implSocket{
		eventHandlers: make(eventHandlers),
		conn:          conn,
		handshake:     newHandshake(conn.Transport().GetRequest()),
	}
	conn.OnClose(func(_ string) {
		if socket.nsp != nil {
			socket.nsp.appendSocket(socket)
			socket.nsp = nil
		}
	})
	conn.OnMessage(func(data []byte) {
		packet, err := parser.Decode(data)
		if err != nil {
			conn.Close()
			return
		}
		switch packet.Type {
		case parser.CONNECT:
			nsp, ok := server.loadNamespace(packet.Namespace)
			if !ok {
				conn.Close()
				if server.logger.err != nil {
					server.logger.err("no such namespace %s\n", packet.Namespace)
				}
				return
			}
			socket.nsp = nsp
			nsp.removeSocket(socket)
			if err := conn.Send(data); err != nil {
				conn.Close()
				if server.logger.err != nil {
					server.logger.err("send connect response failed: %s\n", err)
				}
			}
			break
		case parser.DISCONNECT:
			conn.Close()
			break
		case parser.EVENT:
			// handle income event.
			if model, err := packet.ToModel(); err != nil {
				conn.Close()
				if server.logger.err != nil {
					server.logger.err("%s\n", err)
				}
			} else {
				socket.accept(model.(*parser.MEvent))
			}
			break
		case parser.ACK:
			break
		}
	})
	// add into default namespace.
	socket.nsp = server.Of(defaultNamespace).(*implNamespace)
	socket.nsp.removeSocket(socket)
	if err := conn.Send(initData); err != nil {
		conn.Close()
	}
	return socket
}
