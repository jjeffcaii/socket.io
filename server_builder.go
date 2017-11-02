package sio

import (
	"github.com/jjeffcaii/engine.io"
)

// ServerBuilder can be used to build a new server.
type ServerBuilder struct {
	eio.EngineBuilder
}

// Build returns a new server.
func (p *ServerBuilder) Build() Server {
	return newServer(p.EngineBuilder.Build())
}

// NewBuilder returns a builder for server.
func NewBuilder() *ServerBuilder {
	var builder = &ServerBuilder{
		EngineBuilder: *eio.NewEngineBuilder(),
	}
	builder.SetPath(DefaultPath)
	return builder
}
