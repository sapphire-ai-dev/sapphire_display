package server

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

type displayViewer struct {
	conn      *websocket.Conn
	worldName string
	actorId   int
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

	viewer := &displayViewer{
		worldName: worldName,
		conn:      c,
	}
	registeredWorlds[worldName].viewers = append(registeredWorlds[worldName].viewers, viewer)
	go viewerListener(viewer)
}

func viewerListener(viewer *displayViewer) {
	for {
		_, message, err := viewer.conn.ReadMessage()
		if err != nil {
			fmt.Println("4", err)
			break
		}

		//fmt.Println("viewer", viewer.actorId, "sent", "\""+string(message)+"\" to", viewer.worldName)
		viewerSpeech(viewer.worldName, viewer.actorId, string(message), registeredWorlds[viewer.worldName].conn)
	}
}
