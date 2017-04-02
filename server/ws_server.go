package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var upgrader = websocket.Upgrader{}
var client = &http.Client{}

var port = flag.Int("p", 8888, "websocket server port")
var debug = flag.Bool("d", false, "debug")

func wsHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print(err)
		return
	}
	defer c.Close()

	// XXX: handle disconnects cleanly
	for {
		msgType, msg, err := c.ReadMessage()
		if err != nil {
			log.Print(err)
			break
		}

		if *debug {
			fmt.Printf("Got msg type: %s\n", msgType)
			fmt.Printf("raw msg: %s\n", string(msg))
		}

		// parse message and proxy http request

		if *debug {
			fmt.Printf("got request:\n%s", string(msg))
		}

		buf := bytes.NewBuffer(msg)

		reader := bufio.NewReader(buf)

		req, err := http.ReadRequest(reader)
		if err != nil {
			log.Println(err.Error())
			return
		}

		// See:
		// http://stackoverflow.com/questions/19595860/http-request-requesturi-field-when-making-request-in-go

		// XXX: clean this up
		u, err := url.Parse(fmt.Sprintf(
			"https://%s%s",
			req.Host,
			req.RequestURI,
		))
		if err != nil {
			log.Println(err.Error())
			return
		}

		// this must be cleared out
		req.RequestURI = ""
		req.URL = u

		res, err := client.Do(req)
		if err != nil {
			log.Println(err.Error())
			return
		}

		rawRes, err := httputil.DumpResponse(res, true)
		if err != nil {
			log.Println(err.Error())
			return
		}

		err = c.WriteMessage(msgType, rawRes)
		if err != nil {
			log.Println(err)
		}
	}
}

func main() {
	flag.Parse()
	http.HandleFunc("/ws", wsHandler)
	http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", *port), nil)
}
