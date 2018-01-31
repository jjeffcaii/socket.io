package example

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/jjeffcaii/socket.io"
)

var server sio.Server

func init() {
	flag.Parse()
	server = sio.NewBuilder().Build()
}

func TestEcho(t *testing.T) {
	nsp := server.Of("/")
	nsp.OnConnect(func(socket sio.Socket) {
		socket.On("test", func(msg sio.Message) {
			fmt.Printf("[test] <= %v\n", msg.Any())
			socket.Emit("test", "你好，客户端！")
		})
		socket.OnClose(func(reason string) {
			fmt.Println("socket", socket.ID(), "closed")
		})
		socket.Join("party").To("party").Emit("foo", fmt.Sprintf("msg from %s!", socket.ID()))
	})

	time.AfterFunc(30*time.Second, func() {
		nsp.To("party").Emit("foo", fmt.Sprintf("global test message"))
	})

	http.HandleFunc(sio.DefaultPath, server.Router())
	log.Fatalln(http.ListenAndServe(":3000", nil))
}
