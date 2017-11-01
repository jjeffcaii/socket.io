package parser

import (
	"fmt"
	"testing"

	"encoding/json"
)

func TestDecode(t *testing.T) {
	input := []byte(`42/chat,["say",{"to":{"id":"bar","type":0},"content":{"media":0,"body":"fuck"}}]`)
	packet, err := Decode(input)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("type:", packet.Type)
	fmt.Println("namespace:", packet.Namespace)
	fmt.Println("data:", string(packet.Data))

	/*if m, err2 := packet.ToModel(); err2 != nil {
		fmt.Println(err2)
		return
	} else {
		fmt.Printf("model: %+v", m)
	}*/

	input = []byte("0/chat")

	packet, err = Decode(input)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("type:", packet.Type)
	fmt.Println("namespace:", packet.Namespace)
	fmt.Println("data:", string(packet.Data))
}

func TestJSONARR(t *testing.T) {
	input := []byte(`["say",{"to":{"id":"bar","type":0},"content":{"media":0,"body":"fuck"}}]`)
	vv := make([]interface{}, 0)
	if err := json.Unmarshal(input, &vv); err != nil {
		t.Error(err)
	}
	for _, it := range vv {
		fmt.Printf("vv: %+v\n", it)
	}
}
