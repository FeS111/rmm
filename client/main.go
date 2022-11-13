package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"festech.de/rmm/client/config"
	"festech.de/rmm/client/system"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var interrupt = make(chan os.Signal, 1)
var isInterrupt = false

func main() {
	flag.Parse()
	log.SetFlags(0)

	signal.Notify(interrupt, os.Interrupt)

	config.ReadConfiguration()
	system.GetMacAddress()

	u := url.URL{Scheme: "ws", Host: *addr, Path: fmt.Sprintf("/ws/%s", config.Device.DeviceID), RawQuery: "token=123"}
	connectWebsocket(u.String())
}

func connectWebsocket(url string) {
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		tryReconnect(url)
	}
	fmt.Println("Connected to server")
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				tryReconnect(url)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			isInterrupt = true
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				os.Exit(0)
			}
		}
	}
}

func tryReconnect(url string) {
	fmt.Println("trying to reconnect to server")
	if isInterrupt {
		return
	}
	time.Sleep(time.Second * 5)
	connectWebsocket(url)
}