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

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

type PacketType int8

const (
	CONNECT    PacketType = 0
	DISCONNECT PacketType = 1
	EVENT      PacketType = 2
	ACK        PacketType = 3
	ERROR      PacketType = 4
	EVENTBIN   PacketType = 5
	ACKBIN     PacketType = 6
)

var (
	packetTypes = []byte{
		CONNECT:    '0',
		DISCONNECT: '1',
		EVENT:      '2',
		ACK:        '3',
		ERROR:      '4',
		EVENTBIN:   '5',
		ACKBIN:     '6',
	}
	packetTypesR = map[byte]PacketType{
		'0': CONNECT,
		'1': DISCONNECT,
		'2': EVENT,
		'3': ACK,
		'4': ERROR,
		'5': EVENTBIN,
		'6': ACKBIN,
	}
	errEmptyPacketBytes = errors.New("packet bytes is empty")
)

type Packet struct {
	Namespace string
	Type      PacketType
	ID        uint32
	Data      []byte
}

func (p *Packet) ToModel() (Model, error) {
	switch p.Type {
	default:
		return nil, fmt.Errorf("invalid packet type %d", p.Type)
	case CONNECT:
		mc := &MConnect{
			Namespace: p.Namespace,
		}
		return mc, nil
	case DISCONNECT:
		md := &MDisconnect{
			Namespace: p.Namespace,
		}
		return md, nil
	case EVENT:
		events := make([]json.RawMessage, 0)
		if err := json.Unmarshal(p.Data, &events); err != nil {
			return nil, err
		}
		me := &MEvent{
			ID:        p.ID,
			Namespace: p.Namespace,
		}
		if err := json.Unmarshal(events[0], &me.Event); err != nil {
			return nil, err
		}
		me.Data = make([]interface{}, 0)
		for _, bs := range events[1:] {
			me.Data = append(me.Data, []byte(bs))
		}
		return me, nil
	case ACK:
		ma := &MAck{
			Namespace: p.Namespace,
			ID:        p.ID,
		}
		return ma, nil
	}
}

func NewPacket(packetType PacketType, nsp string, id uint32, data []byte) *Packet {
	switch packetType {
	default:
		panic(fmt.Errorf("invalid packet type: %d", packetType))
	case CONNECT, DISCONNECT, EVENT, ACK, ERROR, EVENTBIN, ACKBIN:
		break
	}
	packet := Packet{
		Namespace: nsp,
		Type:      packetType,
		ID:        id,
		Data:      data,
	}
	return &packet
}

func Decode(input []byte) (*Packet, error) {
	if len(input) < 1 {
		return nil, errEmptyPacketBytes
	}
	packetType, ok := packetTypesR[input[0]]
	if !ok {
		return nil, fmt.Errorf("invalid packet type: %s", string(input[0]))
	}
	var nsp string
	rest := input[1:]
	if length := len(rest); length > 0 && rest[0] == '/' {
		for i, v := range rest {
			if v == ',' {
				nsp = string(rest[:i])
				rest = rest[i+1:]
				break
			}
		}
		if len(nsp) < 1 {
			nsp = string(rest)
			rest = rest[length:]
		}
	}
	return NewPacket(packetType, nsp, 0, rest), nil
}

func Encode(packet *Packet) ([]byte, error) {
	bf := new(bytes.Buffer)
	if err := WriteTo(bf, packet); err != nil {
		return nil, err
	}
	return bf.Bytes(), nil
}

func WriteTo(writer io.Writer, packet *Packet) error {
	switch packet.Type {
	default:
		return fmt.Errorf("invalid packet type: %d", packet.Type)
	case EVENTBIN, ACKBIN:
		return writeToBinary(writer, packet)
	case CONNECT, DISCONNECT, EVENT, ACK, ERROR:
		return writeToString(writer, packet)
	}
}

func writeToString(writer io.Writer, packet *Packet) error {
	if _, err := writer.Write([]byte{packetTypes[packet.Type]}); err != nil {
		return err
	}
	if packet.Namespace != "" && packet.Namespace != "/" {
		if _, err := writer.Write([]byte(packet.Namespace)); err != nil {
			return err
		}
		if _, err := writer.Write([]byte{','}); err != nil {
			return err
		}
	}
	if packet.ID > 0 {
		if _, err := writer.Write([]byte(fmt.Sprintf("%d", packet.ID))); err != nil {
			return err
		}
	}
	if packet.Data != nil && len(packet.Data) > 0 {
		if _, err := writer.Write(packet.Data); err != nil {
			return err
		}
	}
	return nil
}

func writeToBinary(writer io.Writer, packet *Packet) error {
	// TODO: write data as binary
	panic("no implements")
}
