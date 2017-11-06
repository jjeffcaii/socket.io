package sio

import (
	"sync"
)

type implRoom struct {
	name      string
	namespace *implNamespace
	members   map[string]*implSocket
	locker    *sync.RWMutex
}

// join add socket into current room.
// return true if member exists already, else return false.
func (p *implRoom) join(member *implSocket) bool {
	p.locker.Lock()
	defer p.locker.Unlock()
	id := member.ID()
	_, ok := p.members[id]
	if !ok {
		p.members[id] = member
	}
	return ok
}

func (p *implRoom) leave(member *implSocket) bool {
	p.locker.Lock()
	defer p.locker.Unlock()
	id := member.ID()
	_, ok := p.members[id]
	if ok {
		delete(p.members, id)
	}
	return ok
}

func (p *implRoom) has(member *implSocket) bool {
	p.locker.RLock()
	defer p.locker.RUnlock()
	_, ok := p.members[member.ID()]
	return ok
}

func (p *implRoom) broadcast() error {
	// TODO
	return nil
}

func newRoom(nsp *implNamespace, name string) *implRoom {
	r := &implRoom{
		namespace: nsp,
		name:      name,
		locker:    new(sync.RWMutex),
		members:   make(map[string]*implSocket),
	}
	return r
}
