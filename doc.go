// Package sio is a server-side Socket.IO in Golang.
//
// Basic Example:
//
//	var server = sio.NewBuilder().Build()
//	var nsp = server.Of("/")
//	nsp.OnConnect(func(socket sio.Socket) {
//		socket.On("news", func(msg sio.Message) {
//			fmt.Println("[news]:", msg.Any())
//		})
//		socket.Emit("hello", "你好，世界！")
//	})
//	http.HandleFunc(sio.DefaultPath, server.Router())
//	log.Fatalln(http.ListenAndServe(":3000", nil))
//
package sio
