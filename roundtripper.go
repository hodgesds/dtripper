package dtripper

import (
	"fmt"
	"github.com/miekg/dns"
	"net/http"
	"strings"
	"sync"
)

type DNSTripper struct {
	sync.RWMutex
	server     string
	port       int
	serializer Serializer
	client     *dns.Client
}

func NewDNSTripper(
	server string,
	port int,
	net string,
	serializer Serializer,
) *DNSTripper {
	return &DNSTripper{
		client: &dns.Client{
			Net:     net,
			UDPSize: 65535,
		},
		serializer: serializer,
		server:     server,
		port:       port,
	}
}

func (dt *DNSTripper) RoundTrip(
	req *http.Request,
) (*http.Response, error) {
	var err error
	var reply *http.Response

	b, err := dt.serializer.Marshal(req)
	if err != nil {
		return nil, err
	}

	dt.RLock()
	s := dt.server
	p := dt.port
	dt.RUnlock()

	query := fmt.Sprintf("%s.%s", string(b), dns.Fqdn(s))

	m := new(dns.Msg)
	m.SetQuestion(query, dns.TypeTXT)
	// m.RecursionDesired = true

	replyMsg, _, err := dt.client.Exchange(m, fmt.Sprintf("%s:%d", s, p))
	if err != nil {
		return nil, err
	}

	if replyMsg.Rcode != dns.RcodeSuccess {
		return nil, fmt.Errorf("request not successful")
	}

	for _, answer := range replyMsg.Answer {
		ts := strings.Split(answer.String(), "\t")
		raw := strings.Join(ts[4:], "")
		raw = strings.Replace(raw, " ", "", -1)
		raw = strings.Replace(raw, `"`, "", -1)

		reply, err = dt.serializer.Unmarshal([]byte(raw))
		if err != nil {
			continue
		}

		return reply, nil
	}

	return nil, err
}
