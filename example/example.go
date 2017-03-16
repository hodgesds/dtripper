package main

import (
	"flag"
	"github.com/hodgesds/dtripper"
	"io/ioutil"
	"log"
	"net/http"
)

var port = flag.Int("p", 8053, "DNS server port")
var host = flag.String("dns", "localhost", "DNS server")
var net = flag.String("net", "tcp", "DNS network request (udp, tcp)")
var url = flag.String("url", "dns://raw.githubusercontent.com/hodgesds/configs/master/.ctags", "url")

func main() {
	flag.Parse()

	t := &http.Transport{}

	tripper := dtripper.NewDNSTripper(
		*host,
		*port,
		*net,
		&dtripper.DefaultSerializer{},
	)

	t.RegisterProtocol("dns", tripper)

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
