# Socket.IO (WARNING: STILL WORKING!!!!)
Unofficial server-side Socket.IO in Golang.

## Example

``` go
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/jjeffcaii/socket.io"
)

func main() {
	flag.Parse()
	server := sio.NewBuilder().Build()
	nsp := server.Of("/")
	nsp.OnConnect(func(socket sio.Socket) {
		socket.On("news", func(msg sio.Message) {
			fmt.Println("[news]:", msg)
		})
		socket.Emit("hello", "你好，世界！")
	})
	http.HandleFunc(sio.DefaultPath, server.Router())
	log.Fatalln(http.ListenAndServe(":3000", nil))
}

```
