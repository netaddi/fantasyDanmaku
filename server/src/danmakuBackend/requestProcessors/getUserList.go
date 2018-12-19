package requestProcessors

import (
	"net/http"
	"database/sql"
	"io"
	"encoding/json"
	"danmakuBackend/danmakuLib"
)

func GetUserList(w http.ResponseWriter, r * http.Request) {
	danmakuLib.LogHTTPRequest(r)
	w.Header().Set("Access-Control-Allow-Origin", "*")

	config := danmakuLib.GetConfig()
	db, err := sql.Open("mysql", config.DBsource)
	if err != nil {
		println("failed to connect database: ", err.Error())
		io.WriteString(w, "{}")
		db.Close()
		return
	}
	defer db.Close()
	dbQuery := "SELECT reg_code, nickname FROM users WHERE enrolled = TRUE AND permission >= 1;"
	rows, err := db.Query(dbQuery)
	if err != nil {
		println("failed to query database.: ", err.Error())
		io.WriteString(w, "{}")
	} else {
		defer rows.Close()
		userList := make([]map[string]interface{}, 0)
		var regCode string
		var nickname string
		for rows.Next(){
			rows.Scan(&regCode, &nickname)
			userInfo := make(map[string]interface{})
			userInfo["regCode"] = regCode
			userInfo["nickname"] = nickname
			userList = append(userList, userInfo)
		}
		jsonData, _ := json.Marshal(userList)
		io.WriteString(w, string(jsonData))
	}

}