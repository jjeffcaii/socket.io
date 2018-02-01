package example

import (
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/jjeffcaii/socket.io"
)

func TestEcho(t *testing.T) {
	var server = sio.NewBuilder().Build()
	var nsp = server.Of("/")
	nsp.OnConnect(func(socket sio.Socket) {
		socket.On("test", func(msg sio.Message) {
			socket.Emit("test", msg.Any())
		})
		socket.OnClose(func(reason string) {
			fmt.Println("socket", socket.ID(), "closed")
		})
		socket.Join("party").In("party").Emit("test", &struct {
			From string `json:"from"`
			Text string `json:"text"`
		}{socket.ID(), "Hello World!"})
	})
	http.HandleFunc(sio.DefaultPath, server.Router())
	log.Fatalln(http.ListenAndServe(":3000", nil))
}
