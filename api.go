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

// Package sio is a server-side Socket.IO in Golang.
// Basic Example - Echo Server: recieve messages and echo them back.
//
// flag.Parse()
// server := sio.NewBuilder().Build()
// nsp := server.Of("/")
// nsp.OnConnect(func(socket sio.Socket) {
//   socket.On("news", func(msg sio.Message) {
//     fmt.Println("[news]:", msg)
//   })
// 	 socket.Emit("hello", "你好，世界！")
// })
// http.HandleFunc(sio.DefaultPath, server.Router())
// log.Fatalln(http.ListenAndServe(":3000", nil))
package sio

import (
	"net/http"
)

// DefaultPath '/socket.io/' is the default http handler path.
const DefaultPath = "/socket.io/"

// Sockets is an alias of socket map.
type Sockets map[string]Socket
// Message is an alias of event message.
type Message interface{}

// Namespace represents a pool of sockets connected under a given scope identified by a pathname (eg: /chat).
type Namespace interface {
	// ID returns the pathname of namespace.
	ID() string
	// OnConnect register handler for socket created in namespace.
	OnConnect(callback func(socket Socket))
	// GetSockets returns all of sockets in current namespace.
	GetSockets() Sockets
}

// Server is socket.io server.
// You can use ServerBuilder to create it simply.
type Server interface {
	// Of initializes and retrieves the given Namespace by its pathname identifier nsp.
	Of(nsp string) Namespace
	// Router returns a http handler.
	Router() func(http.ResponseWriter, *http.Request)
	// GetSockets returns all sockets of server.
	GetSockets() Sockets
	// Close current server.
	Close()
}

// Socket is the fundamental class for interacting with browser clients.
// A Socket belongs to a certain Namespace (by default /) and uses an underlying Client to communicate.
type Socket interface {
	// ID returns session ID of socket.
	ID() string
	// Namespace returns Namespace of current socket.
	Namespace() Namespace
	// Emit emits an event to the socket identified by the string name.
	Emit(event string, any interface{}) error
	// On register a handler of event identified by the string event.
	On(event string, callback func(msg Message))
	// OnError register a handler of error happend.
	OnError(callback func(error))
	// OnClose register a handler of socket closed.
	OnClose(callback func(reason string))
	// Close close current socket.
	Close()
}
