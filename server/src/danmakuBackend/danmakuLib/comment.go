package danmakuLib

import (
	"encoding/json"
	"fmt"
)

type DanmakuContent struct {
	MessageType string
	Text string
	Color string
	Size int
	CommentType string
}

type FrontendAdminMessage struct {
	MessageType string
	AdminOperation string
	OperationParameter string
}

var DefaultSize = 36
var DefaultType = "normal"


func GetJSON(this interface{}) string {
	b, err := json.Marshal(this)
	if err != nil {
		fmt.Println("jsonize error: ", err.Error())
		return "{}"
	}
	return string(b)
}