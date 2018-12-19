package danmakuLib

import (
	"github.com/gorilla/sessions"
	"net/http"
)

var store = sessions.NewCookieStore([]byte("keyThatDoesNotMatter"))
func GetSession(r *http.Request,w http.ResponseWriter) *sessions.Session {
	session, err := store.Get(r, "sessionID")
	if err != nil {
		panic(err)
	}
	session.Save(r, w)
	return session
}
