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
)

func main() {

	config := lib.GetConfig()
	r := mux.NewRouter()

	// router
	r.HandleFunc("/reg", lib.RegHandler)
	r.HandleFunc("/login", lib.LoginHandler)
	//r.PathPrefix("/").Handler(http.FileServer(http.Dir("./frontend/")))
	r.PathPrefix("/").Handler(
		http.StripPrefix("/", http.FileServer(http.Dir("../../frontend/"))))


	fmt.Println("Started listening at port ", strconv.Itoa(config.Port))

	http.ListenAndServe(":" + strconv.Itoa(config.Port), r)

}
