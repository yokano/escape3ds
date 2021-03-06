package controller

import (
	"net/http"
//	"net/url"
//	"time"
	"appengine"
//	"appengine/taskqueue"
	. "server/model"
	. "server/config"
	. "server/lib"
)

// 引数として渡されたユーザキーに該当するユーザのセッションを開始する。
// ユーザーキーに関連付いたセッションIDを生成して memcache と cookie に保存する。
// ブラウザが提供する cookie に保存されたセッションIDを memcache に保存されたセッションIDを比較して認証する。
func (this *Controller) StartSession(w http.ResponseWriter, r *http.Request, key string) {
	c := appengine.NewContext(r)
	model := NewModel(c)
	sessionId := model.StartSession(key)
	cookie := NewCookie("escape3ds", sessionId, HOSTNAME, "/", 1)
	http.SetCookie(w, cookie)
	
	// セッションタイムアウトの設定（SESSION_TIME_OUT は config.go で設定する）
//	values := url.Values{}
//	values.Add("session_id", sessionId)
//	delay, err := time.ParseDuration(SESSION_TIME_LIMIT)
//	Check(c, err)
//	task := taskqueue.NewPOSTTask("/timeout", values)
//	task.Delay = delay
//	task.Name = sessionId
//	_, err = taskqueue.Add(c, task, "Session")
//	Check(c, err)
//	c.Debugf("セッション開始しました。制限時間を設定しました。")
}

// ブラウザの Cookie に保存されているセッションIDを取得する。
// セッションが存在しない場合は空文字を返す。
func (this *Controller) GetSession(c appengine.Context, r *http.Request) string {
	var result string
	cookie, err := r.Cookie("escape3ds")
	if err == http.ErrNoCookie {
		result = ""
	} else if err != nil {
		result = ""
		c.Errorf(err.Error())
	} else {
		result = cookie.Value
	}
	return result
}

// セッションを終了する。
func (this *Controller) CloseSession(c appengine.Context, sessionId string, w http.ResponseWriter, r *http.Request) {
	model := NewModel(c)
	model.RemoveSession(sessionId)
	this.DeleteCookie(c, w)
}

// セッション ID が保存されたクッキーを削除する
func (this *Controller) DeleteCookie(c appengine.Context, w http.ResponseWriter) {
	cookie := NewCookie("escape3ds", "", HOSTNAME, "/", -1)
	http.SetCookie(w, cookie)
}

// ユーザがログインしているかどうか、クライアントからセッションIDを受け取って調べる。
// ログインが必要な処理の前にかならず実行すること。
// ログインしていなければトップページへ飛ばして空文字を返す。
// 有効なセッションIDを持っていたら対応するユーザIDを返す。
// このとき、最後に操作した時刻を現在として、セッションの有効期限を設定し直す。
func (this *Controller) Session(w http.ResponseWriter, r *http.Request) string {
	c := appengine.NewContext(r)
	sessionId := this.GetSession(c, r)
	if sessionId == "" {
		c.Warningf("セッションIDなしで内部へ入ろうとしました")
		http.Redirect(w, r, "/", 302)
		return ""
	}
	
	model := NewModel(c)
	userKey := model.GetUserKeyFromSession(sessionId)
	if userKey == "" {
		c.Warningf("セッションID: %s に該当するユーザーキが存在しません", sessionId)
		this.DeleteCookie(c, w)
		http.Redirect(w, r, "/", 302)
		return ""
	}
	
//	tasks, err := taskqueue.Lease(c, 1, `taskname`, 30)
//	Check(c, err)
	
	return userKey
}
