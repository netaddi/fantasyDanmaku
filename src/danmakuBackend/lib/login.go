package lib

import (
	"net/http"
	"io"
	"database/sql"
	"github.com/gorilla/sessions"
	"fmt"
)

var store = sessions.NewCookieStore([]byte("keyThatDoesNotMatter"))

func LoginHandler(w http.ResponseWriter, r * http.Request) {
	r.ParseForm()
	LogHTTPRequest(r)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	session, err := store.Get(r, "sessionID")

	for _, v := range r.Form {
		if len(v[0]) < 1 {
			io.WriteString(w, "{\"accepted\": 0, \"errMessage\": \"请填写所有字段\"}")
			return
		}
	}

	config := GetConfig()
	db, err := sql.Open("mysql", config.DBsource)
	if err != nil {
		println("failed to connect database.")
		io.WriteString(w, "{\"accepted\": 0, \"errMessage\": \"failed to connect database.\"}")
		db.Close()
		return
	}

	//stmt, err := db.Prepare("SELECT * FROM users WHERE reg_code=\"?\" AND enrolled=FALSE ")
	rows, err := db.Query("SELECT * FROM users WHERE username=\"" + r.Form.Get("username") +
												"\" AND password=md5(\"" + r.Form.Get("password")+ "\")")
	if err != nil {
		println("failed to query database.: ", err.Error())
		io.WriteString(w, "{\"accepted\": 0, \"errMessage\": \"failed to query database.\"}")
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
		session.Save(r, w)
		io.WriteString(w, "{\"accepted\": 1}")

	} else {
		io.WriteString(w, "{\"accepted\": 0, \"errMessage\": \"用户名或密码错误\"}")
	}

	db.Close()

}