package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

type displayWorld struct {
	state   string
	viewers []*displayViewer
	conn    *websocket.Conn
}

var registeredWorlds = map[string]*displayWorld{}

func worldHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		getWorldHandler(w, r)
	}
}

func getWorldHandler(w http.ResponseWriter, r *http.Request) {
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

	var prevViewers []*displayViewer
	if _, seen := registeredWorlds[worldName]; seen {
		prevViewers = registeredWorlds[worldName].viewers
	} else {
		http.Handle(fmt.Sprintf("/%s/static/", worldName), http.FileServer(http.Dir(worldPath()+"/")))
	}

	registeredWorlds[worldName] = &displayWorld{
		state:   "",
		viewers: prevViewers,
		conn:    c,
	}

	go startListener(c, worldName)
	createViewerActors(worldName, c)
}

func createViewerActors(worldName string, worldConn *websocket.Conn) {
	printErr(worldConn.WriteMessage(websocket.TextMessage, newCreateActorsWorldMsg(worldName)))
}

func viewerSpeech(worldName string, actorId int, speech string, worldConn *websocket.Conn) {
	printErr(worldConn.WriteMessage(websocket.TextMessage, newViewerSpeechMsg(worldName, actorId, speech)))
}

const (
	MsgMethodStart = iota
	MsgMethodUpdateState
	MsgMethodCreateViewerActors
	MsgMethodViewerSpeech
)

// WorldMsg
// data structure used to send message to world
type WorldMsg struct {
	Method        int    `json:"method"`
	ActorCount    int    `json:"actorCount"`
	Speech        string `json:"speech"`
	SpeechActorId int    `json:"speechActorId"`
}

type WorldResp struct {
	Method   int    `json:"method"`
	State    string `json:"state"`
	ActorIds []int  `json:"actorIds"`
	Speech   string `json:"speech"`
}

func newCreateActorsWorldMsg(worldName string) []byte {
	msg := &WorldMsg{
		Method:     MsgMethodCreateViewerActors,
		ActorCount: len(registeredWorlds[worldName].viewers),
	}

	data, err := json.Marshal(msg)
	printErr(err)
	return data
}

func newViewerSpeechMsg(worldName string, actorId int, speech string) []byte {
	msg := &WorldMsg{
		Method:        MsgMethodViewerSpeech,
		SpeechActorId: actorId,
		Speech:        speech,
	}

	data, err := json.Marshal(msg)
	printErr(err)
	return data
}

var processResponseFunc = map[int]func(worldName string, resp *WorldResp){
	MsgMethodUpdateState:        processUpdateStateResponse,
	MsgMethodCreateViewerActors: processCreateViewerActorsResponse,
}

func processResponse(worldName string, resp *WorldResp) {
	if f, seen := processResponseFunc[resp.Method]; seen {
		f(worldName, resp)
	}
}

func processUpdateStateResponse(worldName string, resp *WorldResp) {
	registeredWorlds[worldName].state = resp.State
	for _, viewer := range registeredWorlds[worldName].viewers {
		//fmt.Println("state:", registeredWorlds[worldName].state)
		printErr(viewer.conn.WriteMessage(websocket.TextMessage, []byte(registeredWorlds[worldName].state)))
	}
}

func processCreateViewerActorsResponse(worldName string, resp *WorldResp) {
	for i, viewer := range registeredWorlds[worldName].viewers {
		viewer.actorId = resp.ActorIds[i]
	}
}
