package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func startListener(c *websocket.Conn, name string) {
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			fmt.Println("1", err)
			break
		}

		var resp *WorldResp
		err = json.Unmarshal(message, &resp)
		if err != nil {
			fmt.Println("2", err, string(message))
			break
		}

		processResponse(name, resp)

		//registeredWorlds[name].state = string(message)
		//for _, viewer := range registeredWorlds[name].viewers {
		//	fmt.Println(registeredWorlds[name].state)
		//	printErr(viewer.conn.WriteMessage(websocket.TextMessage, []byte(registeredWorlds[name].state)))
		//}
	}
}
