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
	u, err := url.Parse("http://local/ws")
	if err != nil {
		fmt.Println(err)
	}
	wsHeaders := http.Header{
		"Origin":                   {"http://local"},
		"Sec-WebSocket-Extensions": {"permessage-deflate; client_max_window_bits, x-webkit-deflate-frame"},
	}
	rawConn, err := net.Dial("tcp", *addr)
	if err != nil {
		fmt.Println(err)
	}
	wsConn, resp, err := websocket.NewClient(rawConn, u, wsHeaders, 1024, 1024)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)

	input := make(chan string)
	receive := make(chan string)

	//從server接收訊息
	go func() {
		for {
			_, message, err := wsConn.ReadMessage()
			if err != nil {
				fmt.Println(err)
				break
			}
			receive <- string(message)
		}
	}()

	//stdin 接收中心 and Print
	go func() {
		for {
			fmt.Print("Enter Text:")
			var text string
			fmt.Scanln(&text)
			input <- text
			rs := <-receive
			fmt.Println("Receive:", rs)
		}
	}()

	//訊息寄送處
	for s := range input {
		wsConn.WriteMessage(websocket.TextMessage, []byte(s))
	}
}
