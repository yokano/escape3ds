// リクエスト URL のパターンに該当する処理をする。
// 必要に応じて Model と View を使って処理を進める。
package controller

import (
	"net/http"
	"net/url"
	"appengine"
	"fmt"
	"strings"
	"encoding/json"
	. "server/model"
	. "server/view"
	. "server/config"
	. "server/lib"
)

// コントローラオブジェクト
type Controller struct {
}

// Controller オブジェクトを作成して返す
func NewController() *Controller {
	controller := new(Controller)
	return controller
}

// メソッドを http.HandleFunc() で呼び出し可能な関数型に変換して返す。
// http.HandleFunc() は引数として func(http.ResponseWriter, *http.Request) 型の関数しか渡せない。
// コントローラのメソッドをこの関数型で包むことで http.HandleFunc() から呼び出し可能にする。
func (this *Controller) GetHandler(callback func(this *Controller, w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		callback(this, w, r)
	}
}

// リクエスト URL に合わせて処理を振り分ける
func (this *Controller) Handle() {
	table := make(map[string]func(this *Controller, w http.ResponseWriter, r *http.Request), 21)

	// 通常アクセス
	table["/"]         = (*Controller).Top
	table["/editor"]   = (*Controller).Editor
	table["/gamelist"] = (*Controller).Gamelist
	table["/logout"]   = (*Controller).Logout

	// OAuth 関係
	table["/login_twitter"]     = (*Controller).CallbackTwitter
	table["/callback_twitter"]  = (*Controller).LoginTwitter
	table["/login_facebook"]    = (*Controller).LoginFacebook
	table["/callback_facebook"] = (*Controller).CallbackFacebook

	// アカウント登録関係
	table["/interim_registration"] = (*Controller).InterimRegistration
	table["/registration"]         = (*Controller).Registration

	// Ajax API
	table["/add_user"]    = (*Controller).AddUser
	table["/login"]       = (*Controller).Login
	table["/add_game"]    = (*Controller).AddGame
	table["/delete_game"] = (*Controller).DeleteGame
	table["/upload"]      = (*Controller).Upload
	table["/download"]    = (*Controller).Download
	table["/sync/"]       = (*Controller).SyncHandler

	// 管理者アカウント専用
	table["/debug"]             = (*Controller).Debug
	table["/get_users"]         = (*Controller).GetUsers
	table["/get_interim_users"] = (*Controller).GetInterimUsers

	for url, callback := range table {
		http.HandleFunc(url, this.GetHandler(callback))
	}
}

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

// Twitter ボタンが押された時の処理
func (this *Controller) LoginTwitter(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	oauth := NewOAuth1(c, fmt.Sprintf("http://escape-3ds.appspot.com/callback_twitter"))
	result := oauth.RequestToken("https://api.twitter.com/oauth/request_token")
	oauth.Authenticate(w, r, "https://api.twitter.com/oauth/authenticate", result["oauth_token"])
}

// Twitter でユーザが許可をした時に呼び出されるコールバック関数
func (this *Controller) CallbackTwitter(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	token := r.FormValue("oauth_token")
	verifier := r.FormValue("oauth_verifier")
	
	oauth := NewOAuth1(c, "http://escape-3ds.appspotcom/callback_twitter")
	result := oauth.ExchangeToken(token, verifier, "https://api.twitter.com/oauth/access_token")
	
	view := NewView(c, w)
	model := NewModel(c)
	
	if result["oauth_token"] != "" {
		// ログイン成功
		key := ""
		if model.ExistOAuthUser("Twitter", result["user_id"]) {
			// 既存ユーザ
			params := make(map[string]string, 2)
			params["Type"] = "Twitter"
			params["OAuthId"] = result["user_id"]
			key = model.GetUserKey(params)
		} else {
			// 新規ユーザ
			params := make(map[string]string, 4)
			params["user_type"] = "Twitter"
			params["user_name"] = result["screen_name"]
			params["user_oauth_id"] = result["user_id"]
			params["user_pass"] = ""
			user := model.NewUser(params)
			key = model.AddUser(user)
		}
		if this.GetSession(c, r) == "" {
			this.StartSession(w, r, key)
		}
		view.Gamelist(key)
	} else {
		// ログイン失敗
		view.Login()
		fmt.Fprintf(w, "ログインに失敗しました")
	}
}

// Facebook のログインボタンが押された時の処理
func (this *Controller) LoginFacebook(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	oauth := NewOAuth2(c, FACEBOOK_CLIENT_ID, FACEBOOK_CLIENT_SECRET)
	oauth.RequestAuthorizationCode(w, r, "https://www.facebook.com/dialog/oauth", url.QueryEscape("http://escape-3ds.appspot.com/callback_facebook"))
}

// Facebook へアクセストークンを要求する。
// この関数は Facebook から認証コードをリダイレクトで渡された時に呼ばれる。
// ユーザ情報を格納した map を返す。
func (this *Controller) RequestFacebookToken(w http.ResponseWriter, r *http.Request) map[string]string {
	c := appengine.NewContext(r)
	code := r.FormValue("code")
	oauth := NewOAuth2(c, FACEBOOK_CLIENT_ID, FACEBOOK_CLIENT_SECRET)
	token := oauth.RequestAccessToken(w, r, "https://graph.facebook.com/oauth/access_token", url.QueryEscape("http://escape-3ds.appspot.com/callback_facebook"), code)
	response := oauth.RequestAPI(w, "https://graph.facebook.com/me", token)
	
	// JSON を解析
	type UserInfo struct {
		Id string `json:"id"`
		Name string `json:"name"`
	}
	userInfo := new(UserInfo)
	err := json.Unmarshal(response, userInfo)
	Check(c, err)
	
	result := make(map[string]string, 2)
	result["oauth_id"] = userInfo.Id
	result["name"] = userInfo.Name
	
	return result
}

// Facebook でユーザがアクセスを許可した時に呼び出されるコールバック関数
func (this *Controller) CallbackFacebook(w http.ResponseWriter, r*http.Request) {
	c := appengine.NewContext(r)
	userInfo := this.RequestFacebookToken(w, r)
	
	model := NewModel(c)
	view := NewView(c, w)
	
	key := ""
	if model.ExistOAuthUser("Facebook", userInfo["oauth_id"]) {
		// 既存ユーザ
		params := make(map[string]string, 2)
		params["OAuthId"] = userInfo["oauth_id"]
		params["Type"] = "Facebook"
		key = model.GetUserKey(params)
	} else {
		// 新規ユーザ
		params := make(map[string]string, 4)
		params["user_type"] = "Facebook"
		params["user_name"] = userInfo["name"]
		params["user_oauth_id"] = userInfo["oauth_id"]
		params["user_pass"] = ""
		user := model.NewUser(params)
		key = model.AddUser(user)
	}
	
	if this.GetSession(c, r) == "" {
		this.StartSession(w, r, key)
	}
	view.Gamelist(key)
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

// デバッグツールの表示
func (this *Controller) Debug(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	view := NewView(c, w)
	view.Debug()
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

// ゲーム一覧の表示。ログインした状態でトップページを表示するとここへ飛ぶ。
func (this *Controller) Gamelist(w http.ResponseWriter, r *http.Request) {
	userKey := this.Session(w, r)
	c := appengine.NewContext(r)
	view := NewView(c, w)
	view.Gamelist(userKey)
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

// 仮登録ユーザ一覧の取得。Ajax で呼び出す
func (this *Controller) GetInterimUsers(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	model := NewModel(c)
	interimUsers := model.GetInterimUsers()
	
	// キーと名前だけを返す
	result := make(map[string]string, len(interimUsers))
	for key, val := range interimUsers {
		result[key] = val.Name
	}
	
	bytes, err := json.Marshal(result)
	Check(c, err)
	fmt.Fprintf(w, "%s", bytes)
}

// ユーザ一覧の取得。Ajax で呼び出す
func (this *Controller) GetUsers(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	model := NewModel(c)
	users := model.GetAllUser()
	
	result := make(map[string]string, len(users))
	for key, val := range users {
		result[key] = val.Name
	}
	
	bytes, err := json.Marshal(result)
	Check(c, err)
	fmt.Fprintf(w, "%s", bytes)
}

// 引数として渡されたユーザキーに該当するユーザのセッションを開始する。
// ユーザーキーに関連付いたセッションIDを生成して memcache と cookie に保存する。
// ブラウザが提供する cookie に保存されたセッションIDを memcache に保存されたセッションIDを比較して認証する。
func (this *Controller) StartSession(w http.ResponseWriter, r *http.Request, key string) {
	c := appengine.NewContext(r)
	model := NewModel(c)
	sessionId := model.StartSession(key)
	cookie := NewCookie("escape3ds", sessionId, HOSTNAME, "/", 24)
	http.SetCookie(w, cookie)
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
func (this *Controller) CloseSession(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	sessionId := this.GetSession(c, r)
		
	model := NewModel(c)
	model.RemoveSession(sessionId)
	this.DeleteCookie(c, w)
}

// セッション ID が保存されたクッキーを削除する
func (this *Controller) DeleteCookie(c appengine.Context, w http.ResponseWriter) {
	cookie := NewCookie("escape3ds", "", HOSTNAME, "/", -1)
	http.SetCookie(w, cookie)
}

// 現在の通信にセッションが存在するかどうかを調べる。
// ページ表示時に必ず実行する。
// 持っていなければトップページへ飛ばして空文字を返す。
// 持っていたら対応するユーザIDを返す。
// これによってログインせずに内部ページを表示することを防止する。
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
	
	return userKey
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

