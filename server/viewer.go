package server

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

type displayViewer struct {
	conn    *websocket.Conn
	actorId int
}

func viewerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		getViewerHandler(w, r)
	}
}

func getViewerHandler(w http.ResponseWriter, r *http.Request) {
	worldName := urlParam(r, "name")
	if worldName == "" {
		badRequest(w)
		return
	}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		printErr(err)
		return
	}
	registeredWorlds[worldName].viewers = append(registeredWorlds[worldName].viewers, &displayViewer{
		conn: c,
	})
	go viewerListener(c)
}

func viewerListener(c *websocket.Conn) {
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println("viewer sent", string(message))
	}
}
