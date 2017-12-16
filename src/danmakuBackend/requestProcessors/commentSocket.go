package requestProcessors

import (
	"github.com/gorilla/websocket"
	"net/http"
	"danmakuBackend/danmakuLib"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}


type FrontStruct struct {
	available bool
	conn *websocket.Conn
}


var Frontend = &FrontStruct{false, nil}

func HandleSocket(w http.ResponseWriter, r * http.Request){
	danmakuLib.LogHTTPRequest(r)
	conn, err := upgrader.Upgrade(w, r, nil)
	configs := danmakuLib.GetConfig()
	if err != nil {
		println("websocket error: ", err)
		return
	}
	Frontend.conn = conn
	Frontend.available = true

	conn.WriteMessage(websocket.TextMessage, []byte(configs.WSToken))

	return
}