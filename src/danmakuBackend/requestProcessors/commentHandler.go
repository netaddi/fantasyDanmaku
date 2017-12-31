package requestProcessors

import (
	"net/http"
	"time"
	"database/sql"
	"fmt"
	"danmakuBackend/danmakuLib"
	"html"
	"strings"
)

func CommentHandler(w http.ResponseWriter, r * http.Request){
	r.ParseForm()
	danmakuLib.LogHTTPRequest(r)
	session := danmakuLib.GetSession(r, w)

	lastTime, _ := session.Values["lastTimestamp"].(int64)
	if lastTime == 0 {
		danmakuLib.DenyRequest(w, "请先登录再发送<a href=\\\"login.html\\\">点我登陆</a>")
		return
	}
	timeDifference := time.Now().Unix() - lastTime
	fmt.Print(timeDifference)
	if timeDifference < 2 {
		danmakuLib.DenyRequest(w, "请2秒之后再发送！")
		return
	}
	session.Values["lastTimestamp"] = time.Now().Unix()
	session.Save(r, w)

	//permission, _ := session.Values["permission"].(int)
	//if permission < 0{
	//	danmakuLib.DenyRequest(w, "您的账号因为违规操作被封禁，请联系管理员解封")
	//	return
	//}
	userPermission := danmakuLib.QueryPermission(session.Values["user"].(string))
	if userPermission < 1 {
		danmakuLib.DenyRequest(w, "您的账号因为违规操作被封禁，请联系管理员解封")
		return
	}

	username := session.Values["user"]
	comment := r.Form.Get("text")
	color := r.Form.Get("color")

	// tokenize comment as administrator command
	tokens := strings.Fields(comment)
	if tokens[0] == "//admin" && userPermission > 1 {
		ProcessAdminCommand(tokens)
		danmakuLib.AcceptRequest(w)
		return
	}

	danmakuItem := &danmakuLib.DanmakuContent{
		"danmaku",
		html.EscapeString(comment),
		color,
		danmakuLib.DefaultSize,
		danmakuLib.DefaultType }

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
	if err != nil {
		println("error 101: ", err.Error())
		danmakuLib.DenyRequest(w, "database error. ")
		return
	}
	result, err := stmt.Exec(username, comment, color)
	if err != nil {
		println("error 102: ", err.Error())
		danmakuLib.DenyRequest(w, "database error. ")
		return
	}
	affect, err := result.RowsAffected()
	if err != nil {
		println("error 103: ", err.Error())
		danmakuLib.DenyRequest(w, "database error. ")
		return
	}
	if affect == 1{
		if Frontend.available {
			//Frontend.conn.WriteMessage(websocket.TextMessage, []byte(danmakuItem.GetJSON()))
			Frontend.SendMessage(danmakuLib.GetJSON(danmakuItem))
		}
		danmakuLib.AcceptRequest(w)
	} else {
		danmakuLib.DenyRequest(w, "数据库写入失败")
	}

}