package dtripper

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"net/http"
	"net/http/httputil"
)

type Serializer interface {
	Marshal(*http.Request) ([]byte, error)
	Unmarshal([]byte) (*http.Response, error)
}

type DefaultSerializer struct{}

func (d DefaultSerializer) Marshal(req *http.Request) ([]byte, error) {
	reqBytes, err := httputil.DumpRequest(req, true)
	if err != nil {
		return []byte{}, err
	}

	rawReq := []byte(base64.StdEncoding.EncodeToString(reqBytes))

	// max length of 63 for subdomin
	// XXX: refactor this
	strRawReq := string(rawReq)
	if len(strRawReq) > 63 {
		strRawReq = strRawReq[:63] + "." + strRawReq[63:]
	}

	return []byte(strRawReq), nil
}

func (d DefaultSerializer) Unmarshal(b []byte) (*http.Response, error) {
	replyBytes, err := base64.StdEncoding.DecodeString(string(b))
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(replyBytes)

	reader := bufio.NewReader(buf)

	return http.ReadResponse(reader, nil)
}
