package gremgo

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Client is a container for the gremgo interface TODO: Fix this doc
type Client struct {
	host       string
	connection bool
	reqchan    chan []byte
	reschan    map[string]chan int
	results    map[string]map[string]interface{}
}

// NewClient returns a gremgo client for database interaction
func NewClient(host string) (c Client, err error) {

	// Initializes client

	c.host = "ws://" + host
	c.reqchan = make(chan []byte, 25)
	c.reschan = make(map[string]chan int)
	c.results = make(map[string]map[string]interface{})
	c.connection = true

	// Connect to websocket

	d := websocket.Dialer{}
	ws, _, err := d.Dial(c.host, http.Header{})

	// Write worker
	go func() {
		for c.connection == true {
			select {
			case msg := <-c.reqchan:
				err = ws.WriteMessage(2, msg)
				if err != nil {
					log.Fatal(err)
				}
			default:
			}
		}
	}()

	// Read worker
	go func() {
		for c.connection == true {
			_, msg, err := ws.ReadMessage()
			if err != nil {
				log.Fatal(err)
			}
			if msg != nil {
				go sortResponse(&c, msg) // Send data to sorter
			}
		}
	}()

	return
}