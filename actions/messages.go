package actions

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// MessagesSend default implementation.
func MessagesSend(c buffalo.Context) error {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println(err)
	}

	reader(ws)
	return c.Render(http.StatusOK, r.HTML("messages/send.html"))
}

func reader(conn *websocket.Conn) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message to terminal
		fmt.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
		log.Println("Client Connected")
		err = conn.WriteMessage(1, []byte("Hi Client!"))
		if err != nil {
			log.Println(err)
		}

	}
}

// MessagesRecieve default implementation.
func MessagesRecieve(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.JSON("messages/recieve.html"))
}
