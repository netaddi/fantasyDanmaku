package lib

import (
	"github.com/gorilla/websocket"
	"net/http"
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
	LogHTTPRequest(r)
	conn, err := upgrader.Upgrade(w, r, nil)
	configs := GetConfig()
	if err != nil {
		println("websocket error: ", err)
		return
	}
	Frontend.conn = conn
	Frontend.available = true
	//defer conn.Close()

	//wsWriter, err := conn.NextWriter(websocket.TextMessage)
	//wsWriter.Write([]byte("test string \n"))
	//wsWriter.Close()

	conn.WriteMessage(websocket.TextMessage, []byte(configs.WSToken))

	return
}