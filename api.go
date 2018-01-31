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
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
)

// DefaultPath '/socket.io/' is the default http handler path.
const DefaultPath = "/socket.io/"

// Sockets is an alias of socket map.
type Sockets map[string]Socket

// Message is an alias of event message.
type Message []byte

// Parse parse bytes as json.
func (s Message) Parse(dest interface{}) error {
	return json.Unmarshal(s, dest)
}

// Any parse bytes to anything.
func (s Message) Any() interface{} {
	foo := make([]interface{}, 0)
	bf := &bytes.Buffer{}
	bf.WriteByte('[')
	bf.Write(s)
	bf.WriteByte(']')
	err := json.Unmarshal(bf.Bytes(), &foo)
	if err != nil {
		return nil
	}
	return foo[0]
}

func (s Message) String() string {
	return string(s)
}

// Namespace represents a pool of sockets connected under a given scope identified by a pathname (eg: /chat).
type Namespace interface {
	// ID returns the pathname of namespace.
	ID() string
	// OnConnect register handler for socket created in namespace.
	OnConnect(callback func(socket Socket))
	// GetSockets returns all of sockets in current namespace.
	GetSockets() Sockets
	// To targets a room when emitting.
	To(room string) Emitter
	// In targets a room when emitting.
	In(room string) Emitter
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
}

// Emitter emits event to sockets.
type Emitter interface {
	// Emit emits an event to the socket identified by the string name.
	Emit(event string, first interface{}, others ...interface{}) error
}

// Socket is the fundamental class for interacting with browser clients.
// A Socket belongs to a certain Namespace (by default /) and uses an underlying Client to communicate.
type Socket interface {
	// ID returns session ID of socket.
	ID() string
	// Namespace returns Namespace of current socket.
	Namespace() Namespace
	// Handshake returns Handshake of current socket.
	Handshake() *Handshake
	// Emit emits an event to the socket identified by the string name.
	Emit(event string, first interface{}, others ...interface{}) error
	// On register a handler of event identified by the string event.
	On(event string, callback func(msg Message)) error
	// OnError register a handler of error happend.
	OnError(callback func(error)) error
	// OnClose register a handler of socket closed.
	OnClose(callback func(reason string)) error
	// Close close current socket.
	Close()
	// To return an emitter to all members in room exclude yourself.
	To(room string) Emitter
	// In return an emitter to all members in room.
	In(room string) Emitter
	// Join joins a room.
	Join(room string) Socket
	// Leave leaves a room.
	Leave(room string) Socket
	// Leave leaves all the rooms that we've joined.
	LeaveAll() Socket
	// GetProperties returns custom property map for current socket.
	GetProperties() map[string]interface{}
}

// Handshake is the object used when negociating the handshake.
type Handshake struct {
	Headers http.Header
	Address string
	XDomain bool
	URL     string
	Query   url.Values
}
