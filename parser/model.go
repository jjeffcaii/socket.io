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
package parser

import "encoding/json"

type Model interface {
	ToPacket() (*Packet, error)
}

type MConnect struct {
	Namespace string
}

func (p *MConnect) ToPacket() (*Packet, error) {
	return NewPacket(CONNECT, p.Namespace, 0, nil), nil
}

type MDisconnect struct {
	Namespace string
}

func (p *MDisconnect) ToPacket() (*Packet, error) {
	return NewPacket(DISCONNECT, p.Namespace, 0, nil), nil
}

type MAck struct {
	ID        uint32
	Namespace string
}

func (p *MAck) ToPacket() (*Packet, error) {
	return NewPacket(ACK, p.Namespace, p.ID, nil), nil
}

type MError struct {
	Data interface{}
}

func (p *MError) ToPacket() (*Packet, error) {
	var bs []byte
	switch p.Data.(type) {
	default:
		var err error
		bs, err = json.Marshal(p.Data)
		if err != nil {
			return nil, err
		}
		break
	case error:
		bs = []byte(p.Data.(error).Error())
		break
	case string:
		bs = []byte(p.Data.(string))
		break
	case *string:
		bs = []byte(*(p.Data.(*string)))
	case []byte:
		bs = p.Data.([]byte)
	case *[]byte:
		bs = *(p.Data.(*[]byte))
	}
	return NewPacket(ERROR, "", 0, bs), nil
}

type MEvent struct {
	ID        uint32
	Namespace string
	Event     string
	Data      []interface{}
}

func (p *MEvent) ToPacket() (*Packet, error) {
	foo := []interface{}{p.Event}
	foo = append(foo, p.Data...)
	bs, err := json.Marshal(foo)
	if err != nil {
		return nil, err
	}
	return NewPacket(EVENT, p.Namespace, p.ID, bs), nil
}
