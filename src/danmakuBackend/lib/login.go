package lib

import (
	"net/http"
	//"io"
	"database/sql"
	//"github.com/gorilla/sessions"
	"fmt"
	"time"
)


func LoginHandler(w http.ResponseWriter, r * http.Request) {
	r.ParseForm()
	LogHTTPRequest(r)

	session := getSession(r, w)

	for _, v := range r.Form {
		if len(v[0]) < 1 {
			denyRequest(w, "请填写所有字段")
			//io.WriteString(w, "{\"accepted\": 0, \"errMessage\": \"请填写所有字段\"}")
			return
		}
	}

	config := GetConfig()
	db, err := sql.Open("mysql", config.DBsource)
	if err != nil {
		println("failed to connect database.")
		denyRequest(w, "failed to connect database.")
		//io.WriteString(w, "{\"accepted\": 0, \"errMessage\": \"failed to connect database.\"}")
		db.Close()
		return
	}

	//stmt, err := db.Prepare("SELECT * FROM users WHERE reg_code=\"?\" AND enrolled=FALSE ")
	dbQuery := fmt.Sprintf("SELECT username, permission FROM users WHERE username='%s' AND password=md5('%s');",
									r.Form.Get("username"), r.Form.Get("password"))
	rows, err := db.Query(dbQuery)
	if err != nil {
		println("failed to query database.: ", err.Error())
		denyRequest(w, "failed to query database.")
		//io.WriteString(w, "{\"accepted\": 0, \"errMessage\": \"failed to query database.\"}")
		db.Close()
		return
	}


	if rows.Next() {
		//session.Values["user"] =
		var username string
		var permission int
		rows.Scan(&username, &permission)
		fmt.Println("user login:", username)
		session.Values["user"] = username
		session.Values["permission"] = permission
		session.Values["lastTimestamp"] = time.Now().Unix() - 2
		session.Save(r, w)
		acceptRequest(w)
		println("login:", username, ", permission=", permission)
		//io.WriteString(w, "{\"accepted\": 1}")

	} else {
		denyRequest(w, "用户名或密码错误")
		//io.WriteString(w, "{\"accepted\": 0, \"errMessage\": \"用户名或密码错误\"}")
	}

	db.Close()

}