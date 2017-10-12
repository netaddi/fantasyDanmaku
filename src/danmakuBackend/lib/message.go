package lib

import (
	"net/http"
	"io"
	"fmt"
)

func denyRequest(w http.ResponseWriter, msg string){
	writeStr := fmt.Sprintf("{\"accepted\": 0, \"errMessage\": \"%s\"}", msg)
	io.WriteString(w, writeStr)
	return
}
func acceptRequest(w http.ResponseWriter)  {
	io.WriteString(w, "{\"accepted\": 1}")
}
