package danmakuLib

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func getPostForm( r * http.Request) string {
	if r.Method == "POST" {
		return fmt.Sprint(r.Form)
		//for k, v := range r.Form {
		//	fmt.Println(k, ": ", strings.Join(v, ""))
		//	fmt.Sprint()
		//	formMap[k] = strings.Join(v, "")
		//}
	} else {
		return ""
	}
}

func LogHTTPRequest( r *http.Request) {
	log := []string{
		time.Now().Format("2006-01-02 15:04:05"),
		r.Method,
		r.URL.String(),
		r.RemoteAddr,
		r.Header.Get("X-Forwarded-For"),
		r.Header.Get("User-Agent"),
		getPostForm(r),
	}
	fmt.Println(strings.Join(log, "|--|"))

}
