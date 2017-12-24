package requestProcessors

import (
	"github.com/gorilla/websocket"
	"net/http"
	"danmakuBackend/danmakuLib"
	"time"
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
var connectionKeepRoutineActivated = false

func HandleSocket(w http.ResponseWriter, r * http.Request){
	danmakuLib.LogHTTPRequest(r)


	// disable origin checker.
	// allow non-origin connection.
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	configs := danmakuLib.GetConfig()
	if err != nil {
		println("websocket error: ", err.Error())
		return
	}
	Frontend.conn = conn
	Frontend.available = true

	// infinitely write message to keep connection
	if !connectionKeepRoutineActivated {
		go func() {
			for {
				Frontend.conn.WriteMessage(websocket.TextMessage, []byte(configs.WSToken))
				time.Sleep(time.Second * 5)
			}
		}()
		connectionKeepRoutineActivated = true
	}

	return
}