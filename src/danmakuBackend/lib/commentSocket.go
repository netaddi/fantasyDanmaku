package lib

import (
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

//func sendWebsoketMessage(message string){
//	conn, err := upgrader.Upgrade()
//}

func handleSocket(w http.ResponseWriter, r * http.Request){
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		println("websocket error: ", err)
		return
	}
	return
}