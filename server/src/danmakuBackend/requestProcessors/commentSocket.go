package requestProcessors

import (
	"danmakuBackend/danmakuLib"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

const WSTokenSendIntervalSecond = 10

type FrontStruct struct {
	mutex sync.Mutex
	available bool
	conn *websocket.Conn
}

type ConnectedMessage struct {
	MessageType string
}

var Frontend = &FrontStruct {
	available: false,
	conn: nil }
var connectionKeepRoutineActivated = false

//var lastWSWriteTime int64
//const minimumPingInterval = 5

//func updateLastWSWriteTime() {
//	lastWSWriteTime = time.Now().Unix()
//}
//
//func checkPingSendingCondition() bool {
//	return time.Now().Unix() - lastWSWriteTime > minimumPingInterval
//}

func (front FrontStruct)SendMessage(message string) error {
	front.mutex.Lock()
	err := front.conn.WriteMessage(websocket.TextMessage, []byte(message))
	front.mutex.Unlock()
	return err
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
		connectionKeepRoutineActivated = true
		for {
			connectedMessage := &ConnectedMessage{configs.WSToken}
			messageString := danmakuLib.GetJSON(connectedMessage)
			_ = Frontend.SendMessage(messageString)
			time.Sleep(time.Second * WSTokenSendIntervalSecond)
		}
	}

	return
}