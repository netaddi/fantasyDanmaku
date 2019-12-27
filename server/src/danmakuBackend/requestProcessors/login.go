package requestProcessors

import (
	"net/http"
	"database/sql"
	"fmt"
	"time"
	"danmakuBackend/danmakuLib"
)


func LoginHandler(w http.ResponseWriter, r * http.Request) {
	_ = r.ParseForm()
	danmakuLib.LogHTTPRequest(r)

	session := danmakuLib.GetSession(r, w)

	for _, v := range r.Form {
		if len(v[0]) < 1 {
			danmakuLib.DenyRequest(w, "请填写所有字段")
			return
		}
	}

	config := danmakuLib.GetConfig()
	db, err := sql.Open("mysql", config.DBsource)
	defer db.Close()

	if err != nil {
		println("failed to connect database.")
		danmakuLib.DenyRequest(w, "无法连接数据库。请重试或联系工作人员。")
		return
	}

	//stmt, err := db.Prepare("SELECT * FROM users WHERE reg_code=\"?\" AND enrolled=FALSE ")
	dbQuery := fmt.Sprintf("SELECT reg_code, permission FROM users WHERE reg_code='%s' AND password=md5('%s');",
									r.Form.Get("regCode"), r.Form.Get("password"))
	rows, err := db.Query(dbQuery)
	if err != nil {
		println("failed to query database.: ", err.Error())
		danmakuLib.DenyRequest(w, "failed to query database.")
		return
	}
	defer rows.Close()
	defer db.Close()

	if rows.Next() {
		var regCode string
		var permission int
		_ = rows.Scan(&regCode, &permission)
		session.Values["user"] = regCode
		danmakuLib.SetMemoryPermission(regCode, permission)
		session.Values["lastTimestamp"] = time.Now().Unix() - 2
		_ = session.Save(r, w)
		danmakuLib.AcceptRequest(w)
	} else {
		danmakuLib.DenyRequest(w, "用户名或密码错误")
	}

	//db.Close()

}