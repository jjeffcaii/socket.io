package sio

import (
	"github.com/jjeffcaii/engine.io"
	"github.com/jjeffcaii/socket.io/parser"
)

type implSocket struct {
	rawSocket eio.Socket
}

func (p *implSocket) Namespace() Namespace {
	panic("implement me")
}

func (p *implSocket) Emit(event string, any interface{}) error {
	panic("implement me")
}

func (p *implSocket) On(event string, callback func(msg Message)) {
	panic("implement me")
}

func (p *implSocket) OnError(callback func(error)) {
	panic("implement me")
}

func (p *implSocket) OnClose(callback func()) {
	panic("implement me")
}

func (p *implSocket) Close() {
	panic("implement me")
}

func (p *implSocket) accept(packet *parser.Packet) error {
	panic("implement me")
}

func (p *implSocket) ID() string {
	return p.rawSocket.ID()
}
