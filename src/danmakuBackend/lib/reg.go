package lib

import (
	//"encoding/json"
	//"fmt"
	"net/http"
	"io"
	"database/sql"
	//"database/sql/driver"
	"fmt"
)

func RegHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	LogHTTPRequest(r)
	getSession(r, w)


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
		denyRequest(w, "failed to connect database")
		//io.WriteString(w, "{\"accepted\": 0, \"errMessage\": \"failed to connect database.\"}")
		db.Close()
		return
	}
	//stmt, err := db.Prepare("SELECT * FROM users WHERE reg_code=\"?\" AND enrolled=FALSE ")
	rows, err := db.Query("SELECT * FROM users WHERE reg_code=\"" + r.Form.Get("regCode") + "\" AND enrolled=FALSE")
	//println("SELECT * FROM users WHERE reg_code=\"" + r.Form.Get("regCode") + "\" AND enrolled=FALSE")
	if err != nil {
		println("failed to query database.: ", err.Error())
		denyRequest(w, "failed to query database")
		//io.WriteString(w, "{\"accepted\": 0, \"errMessage\": \"failed to query database.\"}")
		db.Close()
		return
	}


	if rows.Next() {
		repeatQuery := fmt.Sprintf("SELECT * FROM users WHERE username='%s'", r.Form.Get("username"))
		row, err := db.Query(repeatQuery)
		if row.Next() {
			denyRequest(w, "用户名已被占用，换一个试试看")
			return
		}
		repeatQuery = fmt.Sprintf("SELECT * FROM users WHERE nickname='%s'", r.Form.Get("nickname"))
		row, err = db.Query(repeatQuery)
		if row.Next() {
			denyRequest(w, "昵称已被占用，换一个试试看")
			return
		}


		stmt, err := db.Prepare("UPDATE users SET username=?, nickname=?, password=MD5(?), enrolled=TRUE WHERE reg_code=?")
		result, err := stmt.Exec(r.Form.Get("username"), r.Form.Get("nickname"), r.Form.Get("password"), r.Form.Get("regCode"))
		affect, err := result.RowsAffected()
		defer stmt.Close()
		if err != nil {
			denyRequest(w, "failed to write database.")
			io.WriteString(w, "{\"accepted\": 0, \"errMessage\": \"failed to write database.\"}")
		}
		if affect == 1 {
			acceptRequest(w)
			//io.WriteString(w, "{\"accepted\": 1}")?
		} else {
			denyRequest(w, "failed to write database")
			//io.WriteString(w, "{\"accepted\": 0, \"errMessage\": \"failed to write database.\"}")
		}
	} else {
		denyRequest(w, "邀请码不存在或已注册")
		//io.WriteString(w, "{\"accepted\": 0, \"errMessage\": \"邀请码不存在或已注册\"}")
	}
	db.Close()

}
