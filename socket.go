package sio

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/jjeffcaii/engine.io"
	"github.com/jjeffcaii/socket.io/parser"
)

var initData = "0"

type implSocket struct {
	nsp           *implNamespace
	rawSocket     eio.Socket
	eventHandlers map[string][]func(Message)
}

func (p *implSocket) Namespace() Namespace {
	return p.nsp
}

func (p *implSocket) Emit(event string, any interface{}) error {
	evt := &parser.MEvent{
		Namespace: p.nsp.name,
		Event:     event,
		Data:      any,
	}
	packet, err := evt.ToPacket()
	if err != nil {
		return err
	}
	var bs []byte
	bs, err = parser.Encode(packet)
	if err != nil {
		return err
	}
	switch packet.Type {
	default:
		return fmt.Errorf("unsupport emit packet type: %d", packet.Type)
	case parser.EVENT:
		return p.rawSocket.Send(string(bs))
	case parser.EVENTBIN:
		return p.rawSocket.Send(bs)
	}
}

func (p *implSocket) On(event string, callback func(msg Message)) {
	if v, ok := p.eventHandlers[event]; ok {
		vv := append(v, callback)
		p.eventHandlers[event] = vv
	} else {
		p.eventHandlers[event] = []func(Message){callback}
	}
}

func (p *implSocket) OnError(callback func(error)) {
	panic("implement me")
}

func (p *implSocket) OnClose(callback func(reason string)) {
	p.rawSocket.OnClose(callback)
}

func (p *implSocket) Close() {
	p.rawSocket.Close()
}

func (p *implSocket) ID() string {
	return p.rawSocket.ID()
}

func (p *implSocket) accept(evt *parser.MEvent) {
	handlers, ok := p.eventHandlers[evt.Event]
	if !ok {
		glog.Warningln("no such event which name is", evt.Event)
		return
	}
	for _, fn := range handlers {
		fn(evt.Data)
	}
}

func newSocket(server *implServer, rawSocket eio.Socket) *implSocket {
	socket := &implSocket{
		eventHandlers: make(map[string][]func(Message)),
		rawSocket:     rawSocket,
	}
	rawSocket.OnClose(func(_ string) {
		if socket.nsp != nil {
			socket.nsp.leaveSocket(socket)
			socket.nsp = nil
		}
	})
	rawSocket.OnMessage(func(data []byte) {
		packet, err := parser.Decode(data)
		if err != nil {
			rawSocket.Close()
			return
		}
		switch packet.Type {
		case parser.CONNECT:
			nsp, ok := server.loadNamespace(packet.Namespace)
			if !ok {
				rawSocket.Close()
				glog.Errorf("no such namespace %s", packet.Namespace)
				return
			}
			socket.nsp = nsp
			nsp.joinSocket(socket)
			if err := rawSocket.Send(data); err != nil {
				rawSocket.Close()
				glog.Errorln("send connect response failed:", err)
			}
			break
		case parser.DISCONNECT:
			rawSocket.Close()
			break
		case parser.EVENT:
			// handle income event.
			if model, err := packet.ToModel(); err != nil {
				rawSocket.Close()
				glog.Errorln(err)
			} else {
				socket.accept(model.(*parser.MEvent))
			}
			break
		case parser.ACK:
			break
		}
	})
	socket.nsp = server.Of("/").(*implNamespace)
	socket.nsp.joinSocket(socket)
	if err := rawSocket.Send(initData); err != nil {
		rawSocket.Close()
	}
	return socket
}
