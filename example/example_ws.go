package main

import (
	"flag"
	"fmt"
	"github.com/hodgesds/dtripper"
	"io/ioutil"
	"log"
	"net/http"
)

var port = flag.Int("p", 8888, "websocket server port")
var host = flag.String("dns", "localhost", "websocket server")
var url = flag.String("url", "ws://raw.githubusercontent.com/hodgesds/configs/master/.ctags", "url")

func main() {
	flag.Parse()

	t := &http.Transport{}

	tripper, err := dtripper.NewWsTripper(
		fmt.Sprintf("ws://%s:%d/ws", *host, *port),
		&dtripper.RawSerializer{},
	)
	if err != nil {
		log.Fatalln(err.Error())
	}

	t.RegisterProtocol("ws", tripper)

	c := &http.Client{Transport: t}

	res, err := c.Get(*url)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	println(string(body))
}
