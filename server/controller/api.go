package controller

import (
	"net/http"
	"appengine"
	"fmt"
	"strings"
	"encoding/json"
	. "server/model"
	. "server/view"
	. "server/lib"
	. "server/config"
)

// ユーザの追加。Ajax で呼び出す API。
func (this *Controller) AddUser(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	params := make(map[string]string, 5)
	params["user_type"] = r.FormValue("user_type")
	params["user_name"] = r.FormValue("user_name")
	params["user_pass"] = r.FormValue("user_pass")
	params["user_mail"] = r.FormValue("user_mail")
	params["user_oauth_id"] = r.FormValue("user_oauth_id")
	
	model := NewModel(c)
	user := model.NewUser(params)
	model.AddUser(user)
}

// ログイン。成功か失敗か判断してからページを遷移するために Ajax で呼び出す。
// 成功したらリダイレクト先の URL を含む JSON を返す。
// 失敗したらエラーメッセージを含む JSON を返す。
func (this *Controller) Login(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	mail := r.FormValue("mail")
	pass := r.FormValue("pass")
	
	model := NewModel(c)
	key, _ := model.LoginCheck(mail, pass)
	if key != "" {
		// ログイン成功
		sessionId := this.GetSession(c, r)
		if sessionId == "" {
			this.StartSession(w, r, key)
		}
		fmt.Fprintf(w, `{"result":true, "to":"/gamelist"}`)
	} else {
		// ログイン失敗
		fmt.Fprintf(w, `{"result":false, "message":"メールアドレスまたはパスワードが間違っています"}`)
	}
}

// ログアウト。クッキーと memcache に保存されたセッション情報を削除する。
func (this *Controller) Logout(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	cookie, err := r.Cookie("escape3ds")
	Check(c, err)
	sessionId := cookie.Value
	
	model := NewModel(c)
	model.RemoveSession(sessionId)
	this.DeleteCookie(c, w)
	
	http.Redirect(w, r, "/", 302)
}

// Ajax のリクエストパラメータとして渡されたユーザを仮登録データベースに保存して、
// 本登録のためのメールを送信する。
func (this *Controller) InterimRegistration(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	
	name := r.FormValue("name")
	mail := r.FormValue("mail")
	pass := r.FormValue("password")
	
	model := NewModel(c)
	key := model.InterimRegistration(name, mail, pass)
	
	SendMail(c, "infomation@escape-3ds.appspotmail.com", mail, "仮登録完了のお知らせ", fmt.Sprintf(INTERIM_MAIL_BODY, name, key))
	
	view := NewView(c, w)
	view.InterimRegistration()
}

// 仮登録済みのユーザを本登録する。
// 仮登録データベースからユーザを削除してユーザデータベースへ追加する。
func (this *Controller) Registration(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	key := r.FormValue("key")
	
	model := NewModel(c)
	model.Registration(key)
	
	view := NewView(c, w)
	view.Registration()
}

// ゲームの新規追加
func (this *Controller) AddGame(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	gameName := r.FormValue("game_name")
	gameDescription := r.FormValue("game_description")
	
	if gameName == "" {
		fmt.Fprintf(w, `{"result": false}`)
		c.Warningf("空のゲーム名でゲームを作成しようとしました")
	} else if gameDescription == "" {
		fmt.Fprintf(w, `{"result": false}`)
		c.Warningf("ゲーム説明文が空のゲームを作成しようとしました")
	}

	sessionId := this.GetSession(c, r)
	if sessionId == "" {
		fmt.Fprintf(w, `{"result": false}`)
		c.Warningf("セッションIDなしでゲームを作成しようとしました")
	}
	
	model := NewModel(c)
	userKey := model.GetUserKeyFromSession(sessionId)
	params := make(map[string]string, 4)
	params["name"] = gameName
	params["description"] = gameDescription
	params["thumbnail"] = ""
	params["user_key"] = userKey
	game := model.NewGame(params)
	model.AddGame(game, userKey)
	
	fmt.Fprintf(w, `{"result":true, "name":"%s", "description":"%s"}`, gameName, gameDescription)
}

// ゲームの削除。
// ゲームは所有者しか削除することはできない。
func (this *Controller) DeleteGame(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	sessionId := this.GetSession(c, r)
	gameKey := r.FormValue("game_key")
	
	if sessionId == "" {
		fmt.Fprintf(w, `{"result":false}`)
		c.Warningf("セッションIDなしで deleteGame() が呼び出されました")
		return
	} else if gameKey == "" {
		fmt.Fprintf(w, `{"result":false}`)
		c.Warningf("ゲームキー無しで deleteGame() が呼び出されました")
		return
	}
	
	model := NewModel(c)
//	userKey := model.GetUserKeyFromSession(sessionId)
//	game := model.GetGame(gameKey)
//	if game.UserKey != userKey {
//		fmt.Fprintf(w, `{"result":false}`)
//		c.Warningf("ユーザキー: %s が他のユーザのゲームを削除しようとしました", userKey)
//		return
//	}
	
	model.DeleteGame(gameKey)
	fmt.Fprintf(w, `{"result":true}`)
}

// クライアントからアップロードされたファイルを blobstore に保存して blobkey を返す。。
// Ajax で使う。method は POST。
func (this *Controller) Upload(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	file, fileHeader, err := r.FormFile("file")
	Check(c, err)
	
	model := NewModel(c)
	blobKey := model.AddBlob(file, fileHeader)
	
	fmt.Fprintf(w, `{"blobkey":"%s"}`, blobKey)
}

// blobstore から blogkey に関連付いたファイルをクライアントへ渡す。
// Ajax で使う。method は GET。
func (this *Controller) Download(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	
	blobKey := r.FormValue("blobKey")
	if blobKey == "" {
		c.Warningf("blobKey なしで download が実行されました")
		return
	}
	
	model := NewModel(c)
	contentType, bytes := model.GetBlob(blobKey)
	
	header := w.Header()
	header.Add("Content-Type", contentType)
	_, err := w.Write(bytes)
	Check(c, err)
}

// /sync/* にマッチしたら呼び出される。
// エディタで編集したゲーム情報をサーバと同期するための処理。
// URL を解析して更に細かいハンドラへ処理を割り振る。
func (this *Controller) SyncHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	path := strings.Split(r.URL.String(), "/") // "/sync/[kind]/[id]"
	
	// 先頭の "/" も含まれるため
	if len(path) != 4 {
		c.Warningf("パラメータが不足した状態でsyncが実行されました")
		return
	}
	
	switch path[2] {
	case "game":
		this.SyncGame(c, w, r, path[3])
	case "scene":
		this.SyncScene(c, w, r, path[3])
	}
}

// ゲームデータの同期
func (this *Controller) SyncGame(c appengine.Context, w http.ResponseWriter, r *http.Request, gameKey string) {
	switch r.Method {
	case "POST":
	case "GET":
	case "PUT":
		
		body := make([]byte, r.ContentLength)
		r.Body.Read(body)
		
		game := new(Game)
		json.Unmarshal(body, game)
		
		model := NewModel(c)
//		oldGame := model.GetGame(gameKey)
//		if oldGame.UserKey != this.Session(w, r) {
//			c.Warningf("他者のゲームを削除しようとしました　userKey:%s, gameKey:%s", this.Session(w, r), gameKey)
//			return
//		}
		
		model.UpdateGame(gameKey, game)
		fmt.Fprintf(w, `{}`)
	case "DELETE":
		c.Debugf("DELETE GAME")
	}
}

// シーンデータの同期
func (this *Controller) SyncScene(c appengine.Context, w http.ResponseWriter, r *http.Request, sceneId string) {

}