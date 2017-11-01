package parser

import (
	"fmt"
	"testing"
)

func TestDecode(t *testing.T) {
	input := []byte(`2/chat,["say", {to: {id: "bar", type: 0}, content: {media: 0, body: "ffff"}}]`)
	packet, err := Decode(input)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("type:", packet.Type)
	fmt.Println("namespace:", packet.Namespace)
	fmt.Println("data:", string(packet.Data))

	input = []byte("0/chat")

	packet, err = Decode(input)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("type:", packet.Type)
	fmt.Println("namespace:", packet.Namespace)
	fmt.Println("data:", string(packet.Data))

}
