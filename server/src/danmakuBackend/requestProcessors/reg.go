package requestProcessors

import (
	"net/http"
	"database/sql"
	"fmt"
	"danmakuBackend/danmakuLib"
)

func RegHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	danmakuLib.LogHTTPRequest(r)
	danmakuLib.GetSession(r, w)


	for _, v := range r.Form {
		if len(v[0]) < 1 {
			danmakuLib.DenyRequest(w, "请填写所有字段")
			return
		}
	}


	config := danmakuLib.GetConfig()
	db, err := sql.Open("mysql", config.DBsource)
	if err != nil {
		println("failed to connect database.")
		danmakuLib.DenyRequest(w, "failed to connect database")
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM users WHERE reg_code=\"" + r.Form.Get("regCode") + "\" AND enrolled=FALSE")
	if err != nil {
		println("failed to query database.: ", err.Error())
		danmakuLib.DenyRequest(w, "failed to query database")
		//db.Close()
		return
	}
	defer rows.Close()

	if rows.Next() {
		repeatQuery := fmt.Sprintf("SELECT * FROM users WHERE nickname='%s'", r.Form.Get("nickname"))
		row, err := db.Query(repeatQuery)
		if row.Next() {
			danmakuLib.DenyRequest(w, "昵称已被占用，换一个试试看")
			return
		}


		stmt, err := db.Prepare("UPDATE users SET nickname=?, password=MD5(?), enrolled=TRUE WHERE reg_code=?")
		result, err := stmt.Exec(r.Form.Get("nickname"), r.Form.Get("password"), r.Form.Get("regCode"))
		affect, err := result.RowsAffected()
		defer stmt.Close()
		if err != nil {
			danmakuLib.DenyRequest(w, "failed to write database.")
		}
		if affect == 1 {
			danmakuLib.AcceptRequest(w)
		} else {
			danmakuLib.DenyRequest(w, "failed to write database")
		}
	} else {
		danmakuLib.DenyRequest(w, "邀请码不存在或已注册")
	}
	//db.Close()

}
