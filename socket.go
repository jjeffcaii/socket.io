package sio

import (
  "errors"
  "fmt"
  "net/http"

  "github.com/jjeffcaii/engine.io"
  "github.com/jjeffcaii/socket.io/internal/parser"
)

var initData = "0"

type eventHandlers map[string][]func(Message)

func (s eventHandlers) register(event string, handler func(Message)) error {
  if exists, ok := s[event]; ok {
    s[event] = append(exists, handler)
  } else {
    s[event] = []func(Message){handler}
  }
  return nil
}

func (s eventHandlers) exec(event string, msg Message) {
  handlers, ok := s[event]
  if !ok {
    return
  }
  for _, fn := range handlers {
    fn(msg)
  }
}

type implSocket struct {
  nsp           *implNamespace
  conn          eio.Socket
  eventHandlers eventHandlers
  handshake     *Handshake
  properties    map[string]interface{}
}

func (p *implSocket) GetProperties() map[string]interface{} {
  return p.properties
}

func (p *implSocket) To(room string) Emitter {
  socket, ok := p.nsp.getSocket(room)
  if ok {
    return socket
  }
  return p.nsp.ensureRoom(room).broadcast(p.ID())
}

func (p *implSocket) In(room string) Emitter {
  socket, ok := p.nsp.getSocket(room)
  if ok {
    return socket
  }
  return p.nsp.ensureRoom(room).broadcast()
}

func (p *implSocket) Join(room string) Socket {
  p.nsp.ensureRoom(room).join(p)
  return p
}

func (p *implSocket) Leave(room string) Socket {
  found, ok := p.nsp.rooms.Load(room)
  if ok {
    found.(*implRoom).leave(p)
  }
  return p
}

func (p *implSocket) LeaveAll() Socket {
  p.nsp.rooms.Range(func(key, value interface{}) bool {
    value.(*implRoom).leave(p)
    return true
  })
  return p
}

func (p *implSocket) Handshake() *Handshake {
  return p.handshake
}

func (p *implSocket) Namespace() Namespace {
  return p.nsp
}

func (p *implSocket) Emit(event string, first interface{}, others ...interface{}) error {
  if p.nsp == nil {
    return fmt.Errorf("socket %s is closing", p.ID())
  }
  body := []interface{}{first}
  if len(others) > 0 {
    body = append(body, others...)
  }
  evt := &parser.MEvent{
    Namespace: p.nsp.name,
    Event:     event,
    Data:      body,
  }
  packet, err := evt.ToPacket()
  if err != nil {
    return err
  }
  bs, err := parser.Encode(packet)
  if err != nil {
    return err
  }
  switch packet.Type {
  default:
    return fmt.Errorf("unsupport emit packet type: %d", packet.Type)
  case parser.EVENT:
    return p.conn.Send(string(bs))
  case parser.EVENTBIN:
    return p.conn.Send(bs)
  }
}

func (p *implSocket) On(event string, callback func(msg Message)) error {
  return p.eventHandlers.register(event, callback)
}

func (p *implSocket) OnError(callback func(error)) error {
  panic("implement me")
}

func (p *implSocket) OnClose(callback func(reason string)) error {
  p.conn.OnClose(callback)
  return nil
}

func (p *implSocket) Close() {
  p.conn.Close()
}

func (p *implSocket) ID() string {
  return p.conn.ID()
}

func (p *implSocket) accept(evt *parser.MEvent) {
  handlers, ok := p.eventHandlers[evt.Event]
  if !ok {
    if p.nsp.server.logger.warn != nil {
      p.nsp.server.logger.warn("no such event '%s'\n", evt.Event)
    }
    return
  }
  for _, it := range evt.Data {
    for _, fn := range handlers {
      fn(it.([]byte))
    }
  }
}

func newHandshake(req *http.Request) *Handshake {
  if req == nil {
    panic(errors.New("request is nil"))
  }
  return &Handshake{
    Headers: req.Header,
    Query:   req.URL.Query(),
    URL:     req.URL.Path,
    Address: req.RemoteAddr,
  }
}

func handleSocket(server *implServer, conn eio.Socket) *implSocket {
  socket := &implSocket{
    eventHandlers: make(eventHandlers),
    conn:          conn,
    handshake:     newHandshake(conn.Transport().GetRequest()),
    properties:    make(map[string]interface{}),
  }
  conn.OnMessage(func(data []byte) {
    packet, err := parser.Decode(data)
    if err != nil {
      conn.Close()
      return
    }
    switch packet.Type {
    case parser.CONNECT:
      nsp, ok := server.loadNamespace(packet.Namespace)
      if !ok {
        conn.Close()
        if server.logger.err != nil {
          server.logger.err("no such namespace %s\n", packet.Namespace)
        }
        return
      }
      socket.nsp = nsp
      nsp.addSocket(socket)
      // TODO: support custom packet type(binary/text).
      if err := conn.Send(string(data)); err != nil {
        conn.Close()
        if server.logger.err != nil {
          server.logger.err("send connect response failed: %s\n", err)
        }
      }
      break
    case parser.DISCONNECT:
      conn.Close()
      break
    case parser.EVENT:
      // handle income event.
      if model, err := packet.ToModel(); err != nil {
        conn.Close()
        if server.logger.err != nil {
          server.logger.err("%s\n", err)
        }
      } else {
        socket.accept(model.(*parser.MEvent))
      }
      break
    case parser.ACK:
      break
    }
  })
  // add into default namespace.
  socket.nsp = server.Of(defaultNamespace).(*implNamespace)
  socket.nsp.addSocket(socket)
  conn.OnClose(func(_ string) {
    if socket.nsp != nil {
      socket.nsp.removeSocket(socket)
      socket.nsp = nil
    }
  })
  // send init data.
  if err := conn.Send(initData); err != nil {
    conn.Close()
  }
  return socket
}
