// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
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
	}()

	rb := make([]byte, 256)
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
