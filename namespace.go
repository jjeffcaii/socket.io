package sio

import (
	"errors"
	"fmt"
	"sync"
)

const defaultNamespace = "/"

type implNamespace struct {
	server       *implServer
	name         string
	connHandlers []func(Socket)
	sockets      *sync.Map // key=sid,value=socket
	rooms        *sync.Map // key=room_name,value=implRoom
}

func (p *implNamespace) To(roomName string) Emitter {
	return p.In(roomName)
}

func (p *implNamespace) In(roomName string) Emitter {
	socket, ok := p.getSocket(roomName)
	if ok {
		return socket
	}
	return p.ensureRoom(roomName).broadcast()
}

func (p *implNamespace) ID() string {
	return p.name
}

func (p *implNamespace) OnConnect(callback func(socket Socket)) {
	p.connHandlers = append(p.connHandlers, callback)
}

func (p *implNamespace) GetSockets() Sockets {
	ret := make(Sockets)
	p.sockets.Range(func(key, value interface{}) bool {
		ret[key.(string)] = value.(Socket)
		return true
	})
	return ret
}

func (p *implNamespace) getSocket(sid string) (*implSocket, bool) {
	v, ok := p.sockets.Load(sid)
	if !ok {
		return nil, false
	}
	return v.(*implSocket), true
}

func (p *implNamespace) ensureRoom(roomName string) *implRoom {
	var room *implRoom
	if found, ok := p.rooms.Load(roomName); !ok {
		room = newRoom(p, roomName)
		p.rooms.Store(roomName, room)
	} else {
		room = found.(*implRoom)
	}
	return room
}

func (p *implNamespace) removeSocket(socket *implSocket) error {
	p.sockets.Delete(socket.ID())
	p.rooms.Range(func(key, value interface{}) bool {
		room := value.(*implRoom)
		room.members.Delete(socket.ID())
		return true
	})
	return nil
}

func (p *implNamespace) addSocket(socket *implSocket) (err error) {
	defer func() {
		e := recover()
		if e == nil {
			return
		}
		if p.server.logger.err != nil {
			p.server.logger.err("handle socket create failed: %s\n", e)
		}
		switch e.(type) {
		case error:
			err = e.(error)
			break
		case string:
			err = errors.New(e.(string))
			break
		}
	}()
	sid := socket.ID()
	exist, loaded := p.sockets.LoadOrStore(sid, socket)
	if loaded {
		if exist == socket {
			return nil
		}
		return fmt.Errorf("socket %s exists already", sid)
	}
	for _, fn := range p.connHandlers {
		fn(socket)
	}
	return nil
}

func newNamespace(server *implServer, name string) *implNamespace {
	nsp := implNamespace{
		server:       server,
		name:         name,
		sockets:      new(sync.Map),
		rooms:        new(sync.Map),
		connHandlers: make([]func(Socket), 0),
	}
	return &nsp
}
