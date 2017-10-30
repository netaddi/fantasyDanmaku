package lib

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

var DefaultSize = 12
var DefaultType = "normal"

func (this DanmakuContent) getJSON () string{
	b, err := json.Marshal(this)
	if err != nil {
		fmt.Println("jsonize error ! \n")
		return ""
	}
	return string(b)

}
