package requestProcessors

import (
	"danmakuBackend/danmakuLib"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
)

func GetUserList(w http.ResponseWriter, r *http.Request) {
	danmakuLib.LogHTTPRequest(r)
	w.Header().Set("Access-Control-Allow-Origin", "*")

	config := danmakuLib.GetConfig()
	db, _ := sql.Open("mysql", config.DBsource)

	defer db.Close()
	dbQuery := "SELECT reg_code, nickname FROM users WHERE enrolled = TRUE AND permission >= 1;"
	rows, _ := db.Query(dbQuery)
	defer rows.Close()
	userList := make([]map[string]interface{}, 0)
	var regCode string
	var nickname string
	for rows.Next() {
		_ = rows.Scan(&regCode, &nickname)
		userInfo := make(map[string]interface{})
		userInfo["regCode"] = regCode
		userInfo["nickname"] = nickname
		userList = append(userList, userInfo)
	}
	jsonData, _ := json.Marshal(userList)
	_, _ = io.WriteString(w, string(jsonData))
}

func GetRecentCommentList(w http.ResponseWriter, r *http.Request) {
	danmakuLib.LogHTTPRequest(r)
	w.Header().Set("Access-Control-Allow-Origin", "*")

	config := danmakuLib.GetConfig()
	db, _ := sql.Open("mysql", config.DBsource)

	defer db.Close()
	dbQuery := `select content, time, nickname, reg_code
				from comments
  					left join users
    					on comments.user = users.reg_code
				order by comments.time desc
				limit 100;`

	rows, _ := db.Query(dbQuery)
	defer rows.Close()
	userList := make([]map[string]interface{}, 0)
	var content string
	var sendTime string
	var nickname string
	var regCode string
	for rows.Next() {
		_ = rows.Scan(&content, &sendTime, &nickname, &regCode)
		userInfo := make(map[string]interface{})
		userInfo["content"] = content
		userInfo["time"] = sendTime
		userInfo["reg"] = regCode
		userInfo["nickname"] = nickname
		userList = append(userList, userInfo)
	}
	jsonData, _ := json.Marshal(userList)
	_, _ = io.WriteString(w, string(jsonData))

}

func GetUserCommentRanking(w http.ResponseWriter, r *http.Request) {
	danmakuLib.LogHTTPRequest(r)
	w.Header().Set("Access-Control-Allow-Origin", "*")

	config := danmakuLib.GetConfig()
	db, _ := sql.Open("mysql", config.DBsource)

	defer db.Close()
	dbQuery := `select count(nickname), nickname
                from comments
                  left join users
                    on comments.user = users.reg_code
                where nickname != ''
                group by nickname
                order by count(nickname) desc`

	rows, _ := db.Query(dbQuery)
	defer rows.Close()
	userList := make([]map[string]interface{}, 0)
	var count string
	var nickname string
	for rows.Next() {
		_ = rows.Scan(&count, &nickname)
		userInfo := make(map[string]interface{})
		userInfo["count"] = count
		userInfo["nickname"] = nickname
		userList = append(userList, userInfo)
	}
	jsonData, _ := json.Marshal(userList)
	_, _ = io.WriteString(w, string(jsonData))

}
