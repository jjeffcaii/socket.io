package sio

import (
	"net/http"
)

const DefaultPath = "/socket.io/"

type Sockets map[string]Socket
type Message interface{}

type Namespace interface {
	ID() string
	OnConnect(callback func(socket Socket))
	GetSockets() Sockets
}

type Server interface {
	Of(nsp string) Namespace
	Router() func(http.ResponseWriter, *http.Request)
	GetSockets() Sockets
	Close()
}

type Socket interface {
	ID() string
	Namespace() Namespace
	Emit(event string, any interface{}) error
	On(event string, callback func(msg Message))
	OnError(callback func(error))
	OnClose(callback func(reason string))
	Close()
}
