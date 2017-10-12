package main


import (
	"fmt"
	//"os"
	//"encoding/json"
	"danmakuBackend/lib"
	"github.com/gorilla/mux"
	//"github.com/rs/cors"
	//"github.com/gorilla/handlers"
	"net/http"
	"strconv"
	//"encoding/gob"
)

type lastTimestamp int64

func main() {

	config := lib.GetConfig()
	r := mux.NewRouter()

	//gob.Register(lastTimestamp(0))

	// router
	r.HandleFunc("/reg", lib.RegHandler)
	r.HandleFunc("/login", lib.LoginHandler)
	r.HandleFunc("/send", lib.CommentHanbler)

	r.PathPrefix("/").Handler(
		http.StripPrefix("/", http.FileServer(http.Dir("../../frontend/"))))


	fmt.Println("Started listening at port ", strconv.Itoa(config.Port))

	http.ListenAndServe(":" + strconv.Itoa(config.Port), r)

}
