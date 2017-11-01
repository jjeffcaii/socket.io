package example

import (
	"testing"

	"fmt"

	"log"
	"net/http"

	"flag"

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
			fmt.Printf("[test] <= %+v", msg)
			socket.Emit("test", "你好，客户端！")
		})
	})
	http.HandleFunc(sio.DefaultPath, server.Router())
	log.Fatalln(http.ListenAndServe(":3000", nil))
}
