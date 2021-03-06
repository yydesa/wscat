// +build ignore

package main

import (
	"flag"
	"log"
	"strings"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

type errmsg struct {
	where string;
	what error;
};

func fwd(w http.ResponseWriter, r *http.Request, addr string) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	log.Println ("Connection from", r.RemoteAddr,"to",addr)
	ff := r.Header.Get("X-Forwarded-For")
	if (ff != "") {
		log.Println ("Forwarded for", ff)
	}
	ri := r.Header.Get("X-Real-IP")
	if (ri != "") {
		log.Println ("Real IP", ri)
	}

	defer c.Close()

	conn, err := net.Dial("tcp",addr)
	if err != nil {
		log.Println("connect:", addr, err)
		// Just let the conn die...
		return
	}

	defer conn.Close()

	ch := make (chan errmsg, 2)

	go func () {
		rb := make([]byte, 256)
		for {
			n, err := conn.Read(rb)
			if err != nil {
				ch <- errmsg{"read local", err}
				break
			}
			err = c.WriteMessage(websocket.BinaryMessage, rb[0:n])
			if err != nil {
				ch <- errmsg{"write websock", err}
				break
			}
		}
	} ()

	go func () {
		for {
			_, rb, err := c.ReadMessage()
			if err != nil {
				ch <- errmsg{"read websock", err}
				break
			}
			z := 0
			n := len(rb)
			for z < n {
				m, err := conn.Write(rb[z:n])
				if err != nil {
					ch <- errmsg{"write local", err}
					return
				}
				z += m
			}
		}
	} ()

	em := <- ch
	log.Println(em.where, ": ", em.what);
	// Close by defer :-)
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
	some := false
	for _, a := range flag.Args() {
		i := strings.Index(a,"=")
		if (i < 0) {
			log.Fatal ("Bad arg: ", a);
		}
		path := a[0:i]
		targ := a[i+1:len(a)]
		log.Println("url", path, "to", targ)
		some = true
		http.HandleFunc(
			"/"+path,
			func (w http.ResponseWriter, r *http.Request) {
				fwd(w,r,targ)
			})
	}
	if (!some) {
		http.HandleFunc(
			"/webssh",
			func (w http.ResponseWriter, r *http.Request) {
				fwd(w,r,"localhost:22")
			})
	}
	http.HandleFunc("/echo", echo)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
