package sio

import (
	"flag"
	"log"
	"net/http"
	"testing"
)

func init() {
	flag.Parse()
}

func TestAPI(t *testing.T) {
	server := NewBuilder().Build()
	nsp, _ := server.Of("/chat")
	nsp.OnConnect(func(socket Socket) {
		socket.On("friend", func(msg Message) {
			log.Println("rcv from friend:", msg.ToString())
		})
		socket.Emit("friend", "Hello World!")
	})

	http.HandleFunc("/socket.io", server.Router())
	log.Fatalln(http.ListenAndServe(":3000", nil))
}
