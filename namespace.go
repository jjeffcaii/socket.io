package sio

import (
	"sync"

	"fmt"

	"github.com/golang/glog"
)

const defaultNamespace = "/"

type implNamespace struct {
	server       Server
	name         string
	sockets      *sync.Map
	connHandlers []func(Socket)
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
		k := key.(string)
		v := value.(Socket)
		ret[k] = v
		return false
	})
	return ret
}

func (p *implNamespace) leaveSocket(socket *implSocket) error {
	p.sockets.Delete(socket.ID())
	return nil
}

func (p *implNamespace) joinSocket(socket *implSocket) error {
	sid := socket.ID()
	if exist, loaded := p.sockets.LoadOrStore(sid, socket); loaded {
		if exist == socket {
			return nil
		}
		return fmt.Errorf("socket %s exists already", sid)
	}

	c := len(p.connHandlers)
	if c < 1 {
		return nil
	}

	if c == 1 {
		p.connHandlers[0](socket)
		return nil
	}

	wg := new(sync.WaitGroup)
	wg.Add(len(p.connHandlers))
	for _, fn := range p.connHandlers {
		go func() {
			defer func() {
				if e := recover(); e != nil {
					glog.Errorln("handle socket create failed:", e)
				}
				wg.Done()
			}()
			fn(socket)
		}()
	}
	wg.Wait()
	return nil
}

func newNamespace(server Server, name string) *implNamespace {
	nsp := implNamespace{
		server:       server,
		name:         name,
		sockets:      new(sync.Map),
		connHandlers: make([]func(Socket), 0),
	}
	return &nsp
}
