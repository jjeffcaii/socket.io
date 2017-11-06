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
