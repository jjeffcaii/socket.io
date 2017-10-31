package sio

import (
	"encoding/json"
	"net/http"
)

type Sockets map[string]Socket
type Message []byte

func (p Message) Parse(target interface{}) error {
	return json.Unmarshal(p, target)
}

func (p Message) ToString() string {
	return string(p)
}

type Namespace interface {
	OnConnect(callback func(socket Socket))
	GetSockets() Sockets
	Router() func(http.ResponseWriter, *http.Request)
}

type Server interface {
	Of(nsp string) (Namespace, error)
	GetSockets() Sockets
}

type Socket interface {
	ID() string
	Namespace() Namespace
	Emit(event string, any interface{}) error
	On(event string, callback func(msg Message))
	OnError(callback func(error))
	OnClose(callback func())
}
