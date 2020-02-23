package rtmapi

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"golang.org/x/net/websocket"
)

var webSocketServerAddress string
var once sync.Once

func testServer(ws *websocket.Conn) {
	defer ws.Close()
	io.Copy(ws, ws)
}

func startServer() {
	http.Handle("/test", websocket.Handler(testServer))
	server := httptest.NewServer(nil)
	webSocketServerAddress = server.Listener.Addr().String()
}

func TestConnect(t *testing.T) {
	once.Do(startServer)

	url := fmt.Sprintf("ws://%s%s", webSocketServerAddress, "/test")
	client := NewClient()
	conn, err := client.Connect(url)
	if err != nil {
		t.Errorf("webSocket connection error %#v", err)
		return
	}

	msg := []byte("hello world\n")
	if _, err := conn.Write(msg); err != nil {
		t.Errorf("webSocket connection error. %#v", err)
	}
	smallMsg := make([]byte, 8)

	if _, err := conn.Read(smallMsg); err != nil {
		t.Errorf("error on WebSocket payload recieve %#v", err)
	}

	if !bytes.Equal(msg[:len(smallMsg)], smallMsg) {
		t.Errorf("error on recieved message comparison. expected %q got %q", msg[:len(smallMsg)], smallMsg)
	}

	if err := conn.Close(); err != nil {
		t.Errorf("error on WebSocket connection close %#v", err)
	}
}

func TestDecodePayload(t *testing.T) {
	input := []byte("{\"type\": \"message\", \"channel\": \"C2147483705\", \"user\": \"U2147483697\", \"text\": \"Hello, world!\", \"ts\": \"1355517523.000005\", \"edited\": { \"user\": \"U2147483697\", \"ts\": \"1355517536.000001\"}}")
	event, err := DefaultPayloadDecoder(input)
	if err != nil {
		t.Errorf("expected event. got error: %#v", err)
		return
	}

	if event == nil {
		t.Error("expecting event to be returned")
		return
	}

	switch event.(type) {
	case *Message:
		// OK
	default:
		t.Errorf("expected message event. got %#v", event)
	}
}

func TestDecodeReplyPayload(t *testing.T) {
	input := []byte("{\"ok\": true, \"reply_to\": 1, \"ts\": \"1355517523.000005\", \"text\": \"Hello world\"}")
	event, err := DefaultPayloadDecoder(input)

	if err != nil {
		t.Errorf("got error %#v", err)
	}

	if event != nil {
		t.Errorf("expected nil event but got %#v", event)
	}
}

func TestDecodeReplyPayloadWithErrorStatus(t *testing.T) {
	input := []byte("{\"ok\": false, \"reply_to\": 1, \"ts\": \"1355517523.000005\", \"text\": \"Hello world\"}")
	event, err := DefaultPayloadDecoder(input)

	switch e := err.(type) {
	case nil:
		t.Errorf("error MUST be returned. returned event is %#v", event)
	case *ReplyStatusError:
		if *e.Reply.OK {
			t.Error("reply status error is given, but it says OK")
		}
	default:
		t.Errorf("something went wrong. event: %#v, error: %#v", event, err)
	}
}

func TestDecodePayloadWithUnknownFormat(t *testing.T) {
	input := []byte("{\"foo\": \"bar\"}")
	event, err := DefaultPayloadDecoder(input)

	switch err.(type) {
	case nil:
		t.Errorf("error MUST be returned. returned event is... %#v", event)
	case *PayloadError:
		// O.K.
	default:
		t.Errorf("something wrong with the reply payload decode. returened %#v and %#v", event, err)
	}

	if event != nil {
		t.Errorf("expecting nil event to be returned, but was %#v", event)
	}
}