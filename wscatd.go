// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"flag"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

func fwd(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	conn, err := net.Dial("tcp","localhost:22")
	if err != nil {
		log.Println("connect:", err)
		// Just let the conn die...
		return
	}
	go func () {
		rb := make([]byte, 256)
		for {
			n, err := conn.Read(rb)
			if err != nil {
				log.Printf("read local: ", err)
				break
			}
			err = c.WriteMessage(websocket.BinaryMessage, rb[0:n])
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
		c.Close ()
	} ()

	for {
		_, rb, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		z := 0
		n := len(rb)
		for z < n {
			m, err := conn.Write(rb[z:n])
			if err != nil {
				log.Println("conn write:", err)
				return
			}
			z += m
		}
	}
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/webssh", fwd)
	http.HandleFunc("/echo", echo)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
