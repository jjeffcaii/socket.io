package parser

import (
	"fmt"
	"testing"
)

func TestModelToPacket(t *testing.T) {
	msg := map[string]interface{}{"content": "fuck you!", "to": 321}
	m := &MEvent{
		Namespace: "/chat",
		Event:     "live",
		Data:      []interface{}{msg},
	}
	packet, err := m.ToPacket()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("packet: %s", string(packet.Data))
}
