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
