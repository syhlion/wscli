package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/urfave/cli"
)

var (
	name    string
	version string
	cmdRun  = cli.Command{
		Name:    "run",
		Usage:   "wscli run [host] [-d debug]",
		Aliases: []string{"r"},
		Action:  run,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "debug,d",
				Usage: "open debug mode",
			},
		},
	}
)

func main() {
	cli.AppHelpTemplate += "\nWEBSITE:\n\t\thttps://github.com/syhlion/wscli\n\n"
	gusher := cli.NewApp()
	gusher.Name = name
	gusher.Author = "Scott (syhlion)"
	gusher.Usage = "wscli [cmd] [host] [-d]"
	gusher.UsageText = "very simple to use ws connect for test"
	gusher.Version = version
	gusher.Compiled = time.Now()
	gusher.Commands = []cli.Command{
		cmdRun,
	}
	gusher.Run(os.Args)

}

func run(c *cli.Context) (err error) {
	addr := c.Args().Get(0)

	wsHeaders := http.Header{}
	wsConn, _, err := websocket.DefaultDialer.Dial(addr, wsHeaders)
	if err != nil {
		fmt.Fprintf(c.App.Writer, "Usage %s\n err %v\n", c.Command.Usage, err)
		return
	}
	defer wsConn.Close()

	closeSign := make(chan int)
	//從server接收訊息
	go func() {
		defer func() {
			wsConn.Close()
			closeSign <- 1
		}()
		for {
			_, message, err := wsConn.ReadMessage()
			if err != nil {
				log.Printf("read err %v\n", err)
				return
			}
			if c.Bool("debug") {
				log.Printf("receive:%s\n", string(message))
			}
		}
	}()

	go func() {
		buf := bufio.NewScanner(os.Stdin)

		for {
			if !buf.Scan() {
				break
			}
			text := buf.Text()
			if err != nil {
				log.Printf("read stdin err %#v\n", err)
			}
			if err := wsConn.WriteMessage(websocket.TextMessage, []byte(text)); err != nil {
				log.Printf("write er %v\n", err)
				return
			} else {
				log.Println("send scuess!")
			}
		}
	}()
	<-closeSign
	log.Println("close connect")
	return

}
