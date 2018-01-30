package sio

import (
	"fmt"
	"strings"
	"sync"
)

type implRoom struct {
	name      string
	namespace *implNamespace
	members   *sync.Map
	locker    *sync.RWMutex
}

// join add socket into current room.
// return true if member exists already, else return false.
func (p *implRoom) join(member *implSocket) bool {
	_, ok := p.members.Load(member.ID())
	p.members.Store(member.ID(), member)
	return ok
}

func (p *implRoom) leave(member *implSocket) bool {
	p.members.Delete(member.ID())
	return true
}

func (p *implRoom) has(member *implSocket) bool {
	_, ok := p.members.Load(member.ID())
	return ok
}

type broadcastEmitter []*implSocket

func (s broadcastEmitter) Emit(event string, first interface{}, others ...interface{}) error {
	texts := make([]string, 0)
	for _, it := range s {
		if err := it.Emit(event, first, others...); err != nil {
			texts = append(texts, err.Error())
		}
	}
	if len(texts) < 1 {
		return nil
	}
	return fmt.Errorf("broadcast emit error: %s", strings.Join(texts, ";"))
}

func (p *implRoom) broadcast(excludes ...string) broadcastEmitter {
	foo := make([]*implSocket, 0)
	if len(excludes) > 0 {
		exmap := make(map[string]bool)
		for _, it := range excludes {
			exmap[it] = true
		}
		p.members.Range(func(key, value interface{}) bool {
			sid := key.(string)
			if _, ok := exmap[sid]; !ok {
				foo = append(foo, value.(*implSocket))
			}
			return true
		})
	} else {
		p.members.Range(func(key, value interface{}) bool {
			foo = append(foo, value.(*implSocket))
			return true
		})
	}
	return broadcastEmitter(foo)
}

func newRoom(nsp *implNamespace, name string) *implRoom {
	r := &implRoom{
		namespace: nsp,
		name:      name,
		locker:    new(sync.RWMutex),
		members:   new(sync.Map),
	}
	return r
}
