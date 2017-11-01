package parser

import (
	"fmt"
	"testing"
)

func TestModelToPacket(t *testing.T) {
	m := &MEvent{
		Namespace: "/chat",
		Event:     "live",
		Data:      map[string]interface{}{"content": "fuck you!", "to": 321},
	}
	packet, err := m.ToPacket()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("packet: %s", string(packet.Data))
}
