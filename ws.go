package dtripper

import (
	"github.com/gorilla/websocket"
	"net/http"
)

// WsTripper is a RoundTripper interface that uses websocket as the transport
// layer.
type WsTripper struct {
	conn       *websocket.Conn
	serializer Serializer
}

// NewWsTripper returns a *WsTripper
func NewWsTripper(
	url string,
	serializer Serializer,
) (*WsTripper, error) {
	// XXX: check the http response
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}

	return &WsTripper{
		conn:       conn,
		serializer: serializer,
	}, nil
}

// Stop shuts down any running goroutines and stops the websocket connection.
func (ws *WsTripper) Stop() {
	ws.conn.Close()
}

// RoundTrip performs the http.Request and returns a http.Response
func (ws *WsTripper) RoundTrip(
	req *http.Request,
) (*http.Response, error) {
	// marshal the request with the serializer for transport
	b, err := ws.serializer.Marshal(req)
	if err != nil {
		return nil, err
	}

	err = ws.conn.WriteMessage(websocket.TextMessage, b)
	if err != nil {
		return nil, err
	}

	// The tricky part is handling the response... Ideally the websocket
	// consumer should be running in a separate goroutine. This means that
	// some sort of synchronization needs to take place. For now it will be
	// dumb and just read from the connection in the same goroutine.
	// XXX: not goroutine safe
	_, msg, err := ws.conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	return ws.serializer.Unmarshal(msg)
}
