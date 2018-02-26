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
	inputs := make([]interface{}, 0)
	inputs = append(inputs, p.Event)
	inputs = append(inputs, p.Data...)
	bs, err := json.Marshal(inputs)
	if err != nil {
		return nil, err
	}
	return NewPacket(EVENT, p.Namespace, p.ID, bs), nil
}
