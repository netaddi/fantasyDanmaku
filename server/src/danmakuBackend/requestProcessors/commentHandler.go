package requestProcessors

import (
	"danmakuBackend/danmakuLib"
	"database/sql"
	"github.com/gorilla/sessions"
	"html"
	"net/http"
	"strings"
	"time"
)

func checkEligibility(r * http.Request, session * sessions.Session) (bool, string){
	lastTime, _ := session.Values["lastTimestamp"].(int64)

	// check login status
	if lastTime == 0 {
		return false, "请先登录再发送<a href=\\\"login.html\\\">点我登陆</a>"
	}

	userPermission := danmakuLib.QueryPermission(session.Values["user"].(string))

	// check user permission
	if userPermission < 1 {
		return false, "您的账号因为违规操作被封禁，请联系管理员解封"
	}
	if userPermission > 1 {
		// tokenize comment as administrator command
		tokens := strings.Fields(r.Form.Get("text"))
		if tokens[0] == "//admin" && userPermission > 1 {
			ProcessAdminCommand(tokens)
			return false, "Admin command sent."
		}
	}
	// check send time interval
	timeDifference := time.Now().Unix() - lastTime
	if timeDifference < 2 {
		return false, "请2秒之后再发送！"
	}
	session.Values["lastTimestamp"] = time.Now().Unix()
	return true, ""
}

func CommentHandler(w http.ResponseWriter, r * http.Request){
	_ = r.ParseForm()
	danmakuLib.LogHTTPRequest(r)
	session := danmakuLib.GetSession(r, w)

	sendEligible, msg := checkEligibility(r, session)
	if !sendEligible {
		danmakuLib.DenyRequest(w, msg)
		return
	}

	_ = session.Save(r, w)

	username := session.Values["user"]
	comment := r.Form.Get("text")
	color := r.Form.Get("color")

	danmakuItem := &danmakuLib.DanmakuContent{
		MessageType: "danmaku",
		Text:        html.EscapeString(comment),
		Color:       color,
		Size:        danmakuLib.DefaultSize,
		CommentType: danmakuLib.DefaultType,
	}

	config := danmakuLib.GetConfig()
	db, err := sql.Open("mysql", config.DBsource)
	if err != nil {
		println("failed to connect database.")
		danmakuLib.DenyRequest(w, "failed to connect database.")
		db.Close()
		return
	}
	defer db.Close()


	stmt, err := db.Prepare("INSERT INTO comments (user, content, time, color) VALUES (?, ?, now(), ?);")
	defer stmt.Close()
	result, err := stmt.Exec(username, comment, color)
	affect, err := result.RowsAffected()
	if affect == 1 {
		if Frontend.available {
			Frontend.SendMessage(danmakuLib.GetJSON(danmakuItem))
		}
		danmakuLib.AcceptRequest(w)
	} else {
		danmakuLib.DenyRequest(w, "数据库写入失败")
	}

}