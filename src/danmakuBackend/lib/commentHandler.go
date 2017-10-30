package lib

import (
	"net/http"
	//"github.com/gorilla/sessions"
	//"github.com/gorilla/sessions"
	"time"
	"database/sql"
	//"go/token"
	"github.com/gorilla/websocket"
)

func CommentHanbler(w http.ResponseWriter, r * http.Request){
	r.ParseForm()
	LogHTTPRequest(r)
	session := getSession(r, w)



	lastTime, _ := session.Values["lastTimestamp"].(int64)
	if lastTime == 0 {
		denyRequest(w, "请先登录再发送")
		return
	}
	timeDifference := time.Now().Unix() - lastTime
	if timeDifference < 2 {
		denyRequest(w, "请2秒之后再发送！")
		return
	}

	permission, _ := session.Values["permission"].(int)
	if permission < 0{
		denyRequest(w, "您的账号因为违规操作被封禁，请联系管理员解封")
		return
	}



	username := session.Values["user"]
	comment := r.Form.Get("text")
	color := r.Form.Get("color")

	danmakuItem := &DanmakuContent{comment, color, DefaultSize, DefaultType}

	config := GetConfig()
	db, err := sql.Open("mysql", config.DBsource)
	if err != nil {
		println("failed to connect database.")
		denyRequest(w, "failed to connect database.")
		db.Close()
		return
	}

	stmt, err := db.Prepare("INSERT INTO comments (user, content, time, color) VALUES (?, ?, now(), ?);")
	defer stmt.Close()
	if err != nil {
		println("error: ", err.Error())
		denyRequest(w, "database error. ")
	}
	result, err := stmt.Exec(username, comment, color)
	if err != nil {
		println("error: ", err.Error())
		denyRequest(w, "database error. ")
	}
	affect, err := result.RowsAffected()
	if err != nil {
		println("error: ", err.Error())
		denyRequest(w, "database error. ")
	}
	if affect == 1{
		if Frontend.available {
			Frontend.conn.WriteMessage(websocket.TextMessage, []byte(danmakuItem.getJSON()))
		}
		acceptRequest(w)
	} else {
		denyRequest(w, "数据库写入失败")
	}

}