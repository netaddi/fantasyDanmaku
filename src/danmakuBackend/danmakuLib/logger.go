package danmakuLib

import (
	"net/http"
	"fmt"
	"time"
	"strings"
)

func LogHTTPRequest( r *http.Request) {
	fmt.Println("{")
	fmt.Println("time:", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("ip:", r.RemoteAddr)
	fmt.Println("x-forwarded-for:", r.Header.Get("X-Forwarded-For"))
	fmt.Println("method:", r.Method)
	fmt.Println("url:", r.URL.String())
	fmt.Println("UA:", r.Header.Get("User-Agent"))
	fmt.Println("cookies:", r.Header.Get("Cookie"))
	if r.Method == "POST" {
		fmt.Println("form: {")
		//fmt.Println(r.Form)
		for k, v := range r.Form {
			fmt.Println(k, ": ", strings.Join(v, ""))
		}
		fmt.Println("  }")
	}
	fmt.Println("}")
	fmt.Println()
}
