# Socket.IO

[![Build Status](https://travis-ci.org/jjeffcaii/socket.io.svg?branch=master)](https://travis-ci.org/jjeffcaii/socket.io)

Unofficial server-side [Socket.IO](https://socket.io) in Golang.

## Example

``` go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jjeffcaii/socket.io"
)

func main() {
	server := sio.NewBuilder().Build()
	nsp := server.Of("/")
	nsp.OnConnect(func(socket sio.Socket) {
		socket.On("news", func(msg sio.Message) {
			fmt.Println("[news]:", msg)
		})
		socket.Emit("hello", "你好，世界！")
		socket.Join("chats")
		socket.To("chats").Emit("hello", "Hello World!")
	})
	http.HandleFunc(sio.DefaultPath, server.Router())
	log.Fatalln(http.ListenAndServe(":3000", nil))
}

```

## Documents

Please see [https://godoc.org/github.com/jjeffcaii/socket.io](https://godoc.org/github.com/jjeffcaii/socket.io).
