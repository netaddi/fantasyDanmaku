package danmakuLib

import (
	"encoding/json"
	"fmt"
)

type DanmakuContent struct {
	Text string
	Color string
	Size int
	CommentType string
}

type FrontendAdminMessage struct {
	IsAdminCommand bool
	AdminOperation string
	OperationParameter string
}

var DefaultSize = 36
var DefaultType = "normal"

//func (this DanmakuContent) GetJSON () string{
//	b, err := json.Marshal(this)
//	if err != nil {
//		fmt.Println("jsonize error ! \n")
//		return ""
//	}
//	return string(b)
//}

func GetJSON(this interface{}) string {
	b, err := json.Marshal(this)
	if err != nil {
		fmt.Println("jsonize error: ", err.Error())
		return "{}"
	}
	return string(b)
}