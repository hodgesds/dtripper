package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/miekg/dns"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
)

func parseQuery(m *dns.Msg) {
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeTXT:
			log.Printf("Query for %s\n", q.Name)

			client := &http.Client{}
			// res, err := c.Get("https://www.google.com/search?q=dns")

			subs := strings.Split(q.Name, ".")
			if len(subs) <= 1 {
				println("no subdomains")
				return
			}

			data := strings.Join(subs[:len(subs)-2], "")

			decoded, err := base64.StdEncoding.DecodeString(data)
			if err != nil {
				println(err.Error())
				return
			}

			fmt.Printf("got request:\n%s", string(decoded))

			buf := bytes.NewBuffer(decoded)

			reader := bufio.NewReader(buf)

			req, err := http.ReadRequest(reader)
			if err != nil {
				println(err.Error())
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
				println(err.Error())
				return
			}

			// this must be cleared out
			req.RequestURI = ""
			req.URL = u

			res, err := client.Do(req)
			if err != nil {
				println("client error")
				println(err.Error())
				return
			}

			rawRes, err := httputil.DumpResponse(res, true)
			if err != nil {
				println(err.Error())
				return
			}

			encoded := base64.StdEncoding.EncodeToString(rawRes)

			resData := []string{}

			processing := true
			for {
				if len(encoded) < 255 {
					resData = append(resData, encoded)
					processing = false
				} else {
					resData = append(resData, `"`+encoded[:255]+`"`)
					encoded = encoded[255:]
				}
				if !processing {
					break
				}
			}

			resStr := strings.Join(resData, " ")

			rr, err := dns.NewRR(fmt.Sprintf(`%s 1 IN TXT %s`, q.Name, resStr))
			if err != nil {
				println(err.Error())
				return
			}
			m.Answer = append(m.Answer, rr)
			fmt.Printf("%+v\n", rr)
		}
	}
}

func handleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = true

	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(m)
	}

	w.WriteMsg(m)
}

var port = flag.Int("p", 8053, "port")
var net = flag.String("net", "tcp", "network (udp or tcp)")

func main() {
	flag.Parse()
	dns.HandleFunc(".", handleDnsRequest)

	server := &dns.Server{
		Addr:    ":" + strconv.Itoa(*port),
		Net:     *net,
		UDPSize: 65535,
	}
	log.Printf("Starting at %d\n", *port)

	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}
}
