package main


import (
	"fmt"
	"danmakuBackend/danmakuLib"
	"danmakuBackend/requestProcessors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	_ "github.com/go-sql-driver/mysql"
)

type lastTimestamp int64

func main() {

	config := danmakuLib.GetConfig()
	danmakuLib.InitPermissionMap()

	r := mux.NewRouter()

	// router
	r.HandleFunc("/reg", requestProcessors.RegHandler)
	r.HandleFunc("/login", requestProcessors.LoginHandler)
	r.HandleFunc("/send", requestProcessors.CommentHandler)
	r.HandleFunc("/ws", requestProcessors.HandleSocket)
	r.HandleFunc("/getUserList", requestProcessors.GetUserList)
	r.HandleFunc("/answer", requestProcessors.ProcessAnswering)
	r.HandleFunc("/getQuestionResult", requestProcessors.GetQuestionResult)

	r.PathPrefix("/").Handler(
		http.StripPrefix("/", http.FileServer(http.Dir("../../frontend/"))))


	fmt.Println("Started listening at port ", strconv.Itoa(config.Port))

	http.ListenAndServe(":" + strconv.Itoa(config.Port), r)

}
