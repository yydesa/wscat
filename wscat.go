// +build ignore

package main

import (
	"flag"
	"log"
	"net/url"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var path = flag.String("path", "webssh", "url path base")
var suffix = flag.String("suffix", "", "url suffix")
var theproxy = flag.String("proxy", "", "proxy server")
var plain = flag.Bool("insecure", false, "no ssl")

func main() {
	flag.Parse()
	log.SetFlags(0)

	//interrupt := make(chan os.Signal, 1)
	//signal.Notify(interrupt, os.Interrupt)

	scheme := "wss";
	if (*plain) {
		scheme = "ws";
	}
	u := url.URL{Scheme: scheme, Host: *addr, Path: "/"+*path+*suffix}
	log.Printf("[%s]", u.String())

	dl := websocket.DefaultDialer
	if (*theproxy != "") {
		u, err := url.Parse(*theproxy)
		if err != nil {
			log.Fatal("proxy format:", err)
		}
		dl = &websocket.Dialer{ Proxy: http.ProxyURL(u) }
	}

	c, _, err := dl.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer c.Close()
		defer close(done)
		for {
			_, rb, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			z := 0
			n := len(rb)
			for z < n {
				m, err := os.Stdout.Write(rb[z:n])
				if err != nil {
					panic("failed to write stdout: " + err.Error())
				}
				z += m
			}
		}
	} ()

	rb := make([]byte, 8200)
	for {
		n, err := os.Stdin.Read(rb)
		if err != nil {
			panic("failed to read stdin: " + err.Error())
		}
		err = c.WriteMessage(websocket.BinaryMessage, rb[0:n])
		if err != nil {
			log.Println("write:", err)
			return
		}
	}
}
