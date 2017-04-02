package main

import (
	"flag"
	"fmt"
	"github.com/hodgesds/dtripper"
	"io/ioutil"
	"log"
	"net/http"
)

var dnsPort = flag.Int("dp", 8053, "DNS server port")
var wsPort = flag.Int("wp", 8888, "websocket server port")
var host = flag.String("dns", "localhost", "tunnel server")
var net = flag.String("net", "tcp", "DNS network request (udp, tcp)")
var url = flag.String("url", "dns://raw.githubusercontent.com/hodgesds/configs/master/.ctags", "url")

func main() {
	flag.Parse()

	t := &http.Transport{}

	dnsTripper := dtripper.NewDNSTripper(
		*host,
		*dnsPort,
		*net,
		&dtripper.DefaultSerializer{},
	)
	wsTripper, err := dtripper.NewWsTripper(
		fmt.Sprintf("ws://%s:%d/ws", *host, *wsPort),
		&dtripper.RawSerializer{},
	)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// register the DNS and websocket protocols to be used with the http.Client
	t.RegisterProtocol("dns", dnsTripper)
	t.RegisterProtocol("ws", wsTripper)

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
