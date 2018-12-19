package danmakuLib

import (
	"net/http"
	"io"
	"fmt"
)

func DenyRequest(w http.ResponseWriter, msg string) {
	writeStr := fmt.Sprintf("{\"accepted\": 0, \"errMessage\": \"%s\"}", msg)
	io.WriteString(w, writeStr)
	return
}
func AcceptRequest(w http.ResponseWriter) {
	io.WriteString(w, "{\"accepted\": 1}")
}

