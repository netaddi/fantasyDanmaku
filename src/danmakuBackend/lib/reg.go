package lib

import (
	//"encoding/json"
	//"fmt"
	"net/http"
	"io"
	"database/sql"
	//"database/sql/driver"
)

func RegHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	LogHTTPRequest(r)


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
	rows, err := db.Query("SELECT * FROM users WHERE reg_code=\"" + r.Form.Get("regCode") + "\" AND enrolled=FALSE")
	println("SELECT * FROM users WHERE reg_code=\"" + r.Form.Get("regCode") + "\" AND enrolled=FALSE")
	if err != nil {
		println("failed to query database.: ", err.Error())
		io.WriteString(w, "{\"accepted\": 0, \"errMessage\": \"failed to query database.\"}")
		db.Close()
		return
	}


	if rows.Next() {
		stmt, err := db.Prepare("UPDATE users SET username=?, nickname=?, password=MD5(?), enrolled=TRUE WHERE reg_code=?")
		result, err := stmt.Exec(r.Form.Get("username"), r.Form.Get("nickname"), r.Form.Get("password"), r.Form.Get("regCode"))
		affect, err := result.RowsAffected()
		if err != nil {
			io.WriteString(w, "{\"accepted\": 0, \"errMessage\": \"failed to write database.\"}")
		}
		if affect == 1 {
			io.WriteString(w, "{\"accepted\": 1}")
		} else {
			io.WriteString(w, "{\"accepted\": 0, \"errMessage\": \"failed to write database.\"}")
		}
	} else {
		io.WriteString(w, "{\"accepted\": 0, \"errMessage\": \"邀请码不存在或已注册\"}")
	}
	db.Close()

}
