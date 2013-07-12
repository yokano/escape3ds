package controller

import (
	"net/http"
	"appengine"
	. "server/view"
)

// ログインページの表示。セッションが既に開始されている場合はゲームリストを表示する。
func (this *Controller) Top(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	view := NewView(c, w)
	sessionId := this.GetSession(c, r)
	
	if sessionId != "" {
		http.Redirect(w, r, "/gamelist", 302)
	} else {
		view.Login()
	}
}

// エディタの表示
func (this *Controller) Editor(w http.ResponseWriter, r *http.Request) {
//	userKey := this.Session(w, r)
	c := appengine.NewContext(r)
//	model := NewModel(c)

	gameKey := r.FormValue("game_key")
	if gameKey == "" {
		c.Warningf("ゲームキー無しでゲームを編集しようとしました")
		http.Redirect(w, r, "/", 302)
	}
	
//	game := model.GetGame(gameKey)
//	if game.UserKey != userKey {
//		c.Warningf("ユーザキー: %s が他人のゲーム: %s を編集しようとしました", userKey, gameKey)
//		http.Redirect(w, r, "/gamelist", 302)
//	}
	
	view := NewView(c, w)
	view.Editor(gameKey)
}

// ゲーム一覧の表示。ログインした状態でトップページを表示するとここへ飛ぶ。
func (this *Controller) Gamelist(w http.ResponseWriter, r *http.Request) {
	userKey := this.Session(w, r)
	c := appengine.NewContext(r)
	view := NewView(c, w)
	view.Gamelist(userKey)
}