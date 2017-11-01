package sio

import (
	"github.com/jjeffcaii/engine.io"
)

type ServerBuilder struct {
}

func (p *ServerBuilder) Build() Server {
	eng := eio.NewEngineBuilder().SetPath(DefaultPath).Build()
	return newServer(eng)
}

func NewBuilder() *ServerBuilder {
	return new(ServerBuilder)
}
