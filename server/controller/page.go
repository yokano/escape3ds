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
	this.Session(w, r)
	c := appengine.NewContext(r)

	gameKey := r.FormValue("game_key")
	if gameKey == "" {
		c.Warningf("ゲームキー無しでゲームを編集しようとしました")
		http.Redirect(w, r, "/", 302)
	}
	
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

// テストプレイ
func (this *Controller) Runtime(w http.ResponseWriter, r *http.Request) {
//	this.Session(w, r)
	gameKey := r.FormValue("game_key")
	c := appengine.NewContext(r)
	view := NewView(c, w)
	view.Runtime(gameKey)
}

// イベントエディタ（シーン内のイベント用：クリックされたら何かする）
func (this *Controller) EventEditor(w http.ResponseWriter, r *http.Request) {
	eventKey := r.FormValue("event_key")
	gameKey := r.FormValue("game_key")
	c := appengine.NewContext(r)
	view := NewView(c, w)
	view.EventEditor(gameKey, eventKey)
}

// イベントエディタ（シーン開始時のイベント用：シーンが表示されたら何かする）
func (this *Controller) EnterEventEditor(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	gameKey := r.FormValue("game_key")
	sceneKey := r.FormValue("scene_key")
	view := NewView(c, w)
	view.EnterEventEditor(gameKey, sceneKey)
}

// イベントエディタ（シーン終了時のイベント用：別のシーンへ移動する直前に実行される）
func (this *Controller) LeaveEventEditor(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	gameKey := r.FormValue("game_key")
	sceneKey := r.FormValue("scene_key")
	view := NewView(c, w)
	view.LeaveEventEditor(gameKey, sceneKey)
}