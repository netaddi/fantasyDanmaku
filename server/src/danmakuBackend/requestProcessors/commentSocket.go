package requestProcessors

import (
	"github.com/gorilla/websocket"
	"net/http"
	"danmakuBackend/danmakuLib"
	"time"
)

const WSTokenSendIntervalSecond = 10

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type FrontStruct struct {
	available bool
	conn *websocket.Conn
}

type ConnectedMessage struct {
	MessageType string
}

var Frontend = &FrontStruct{false, nil}
var connectionKeepRoutineActivated = false

var lastWSWriteTime int64
const minimumPingInterval = 5

func updateLastWSWriteTime() {
	lastWSWriteTime = time.Now().Unix()
}

func checkPingSendingCondition() bool {
	return time.Now().Unix() - lastWSWriteTime > minimumPingInterval
}

func (front FrontStruct)SendMessage(messsage string) {
	front.conn.WriteMessage(websocket.TextMessage, []byte(messsage))
}

func HandleSocket(w http.ResponseWriter, r * http.Request){
	danmakuLib.LogHTTPRequest(r)


	// disable origin checker.
	// allow non-origin connection.
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
			},
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
		//go func() {
			for {
				//if checkPingSendingCondition() {
					connectedMessage := &ConnectedMessage{configs.WSToken}
					messageString := danmakuLib.GetJSON(connectedMessage)
					Frontend.conn.WriteMessage(websocket.TextMessage, []byte(messageString))
					time.Sleep(time.Second * WSTokenSendIntervalSecond)
				//}
			}
		//}()
		connectionKeepRoutineActivated = true
	}

	return
}