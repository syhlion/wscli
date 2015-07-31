package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"net/url"
)

var addr = flag.String("addr", "127.0.0.1:8001", "ws address")

func main() {
	flag.Parse()
	u, err := url.Parse("http://local/ws")
	if err != nil {
		fmt.Println(err)
	}
	wsHeaders := http.Header{
		"Origin":                   {"http://local"},
		"Sec-WebSocket-Extensions": {"permessage-deflate; client_max_window_bits, x-webkit-deflate-frame"},
	}
	println(*addr)
	rawConn, err := net.Dial("tcp", *addr)
	if err != nil {
		fmt.Println(err)
	}
	wsConn, resp, err := websocket.NewClient(rawConn, u, wsHeaders, 1024, 1024)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)

	//從server接收訊息
	go func() {
		for {
			_, message, err := wsConn.ReadMessage()
			if err != nil {
				fmt.Println(err)
				break
			}
			fmt.Println("Receive:", string(message))
		}
	}()

	for {
		var text string
		fmt.Scanln(&text)
		fmt.Println("Send:", text)
		if err := wsConn.WriteMessage(websocket.TextMessage, []byte(text)); err != nil {
			fmt.Println("Please ReConnect")
			break
		}
	}

}
