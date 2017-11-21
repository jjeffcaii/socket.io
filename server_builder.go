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
	"log"
	"net/http"

	"github.com/jjeffcaii/engine.io"
)

// ServerBuilder can be used to build a new server.
type ServerBuilder struct {
	linfo         *log.Logger
	lwarn         *log.Logger
	lerr          *log.Logger
	engineBuilder *eio.EngineBuilder
}

// Build returns a new server.
func (p *ServerBuilder) Build() Server {
	return newServer(p.engineBuilder.Build(), p.linfo, p.lwarn, p.lerr)
}

// SetLoggerInfo set logger for INFO
func (p *ServerBuilder) SetLoggerInfo(logger *log.Logger) *ServerBuilder {
	p.engineBuilder.SetLoggerInfo(logger)
	p.linfo = logger
	return p
}

// SetLoggerWarn set logger for WARN
func (p *ServerBuilder) SetLoggerWarn(logger *log.Logger) *ServerBuilder {
	p.engineBuilder.SetLoggerWarn(logger)
	p.lwarn = logger
	return p
}

// SetLoggerError set logger for ERROR
func (p *ServerBuilder) SetLoggerError(logger *log.Logger) *ServerBuilder {
	p.engineBuilder.SetLoggerError(logger)
	p.lerr = logger
	return p
}

// SetTransports define transport types allow.
func (p *ServerBuilder) SetTransports(transports ...eio.TransportType) *ServerBuilder {
	p.engineBuilder.SetTransports(transports...)
	return p
}

// SetGenerateID define the method of creating SocketID.
func (p *ServerBuilder) SetGenerateID(gen func(uint32) string) *ServerBuilder {
	p.engineBuilder.SetGenerateID(gen)
	return p
}

// SetPath define the http router path for Engine.
func (p *ServerBuilder) SetPath(path string) *ServerBuilder {
	p.engineBuilder.SetPath(path)
	return p
}

// SetAllowRequest set a function that receives a given request, and can decide whether to continue or not.
func (p *ServerBuilder) SetAllowRequest(validator func(*http.Request) error) *ServerBuilder {
	p.engineBuilder.SetAllowRequest(validator)
	return p
}

// SetCookie can control enable/disable of cookie.
func (p *ServerBuilder) SetCookie(enable bool) *ServerBuilder {
	p.engineBuilder.SetCookie(enable)
	return p
}

// SetCookiePath define the path of cookie.
func (p *ServerBuilder) SetCookiePath(path string) *ServerBuilder {
	p.engineBuilder.SetCookiePath(path)
	return p
}

// SetCookieHTTPOnly if set true HttpOnly io cookie cannot be accessed by client-side APIs,
// such as JavaScript. (true) This option has no effect
// if cookie or cookiePath is set to false.
func (p *ServerBuilder) SetCookieHTTPOnly(httpOnly bool) *ServerBuilder {
	p.engineBuilder.SetCookieHTTPOnly(httpOnly)
	return p
}

// SetAllowUpgrades define whether to allow transport upgrades. (default allow upgrades)
func (p *ServerBuilder) SetAllowUpgrades(enable bool) *ServerBuilder {
	p.engineBuilder.SetAllowUpgrades(enable)
	return p
}

// SetPingInterval define ping time interval in millseconds for client.
func (p *ServerBuilder) SetPingInterval(interval uint32) *ServerBuilder {
	p.engineBuilder.SetPingInterval(interval)
	return p
}

// SetPingTimeout define ping timeout in millseconds for client.
func (p *ServerBuilder) SetPingTimeout(timeout uint32) *ServerBuilder {
	p.engineBuilder.SetPingTimeout(timeout)
	return p
}

// NewBuilder returns a builder for server.
func NewBuilder() *ServerBuilder {
	var builder = &ServerBuilder{
		engineBuilder: eio.NewEngineBuilder(),
	}
	builder.SetPath(DefaultPath)
	return builder
}
