package controller

import (
	"net/http"
	"appengine"
	"appengine/taskqueue"
	"fmt"
	"strings"
	"time"
	"net/url"
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

// ユーザの削除
func (this *Controller) DeleteUser(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	userKey := r.FormValue("user_key")
	model := NewModel(c)
	model.DeleteUser(userKey)
	this.Logout(w, r)
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

// ゲストとしてログインする
func (this *Controller) GuestLogin(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	model := NewModel(c)
	
	// ゲストユーザ追加
	user := model.NewUser(map[string]string{
		"user_type": "guest",
		"user_name": "ゲスト",
		"user_pass": "",
		"user_mail": "",
		"user_oauth_type": "",
	})
	userKey := model.AddUser(user)
	this.StartSession(w, r, userKey)
	
	// 24時間後にゲストユーザを削除
	values := url.Values{}
	values.Set("user_key", userKey)
	delay, err := time.ParseDuration("24h")
	Check(c, err)

	task := taskqueue.NewPOSTTask("/bye_guest", values)
	task.Delay = delay
	taskqueue.Add(c, task, "default")
	
	view := NewView(c, w)
	view.Gamelist(userKey)
}

// ログアウト。クッキーと memcache に保存されたセッション情報を削除する
// セッションタイムアウト用のタスクも削除する
func (this *Controller) Logout(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	cookie, err := r.Cookie("escape3ds")
	Check(c, err)
	sessionId := cookie.Value
	
	model := NewModel(c)
	model.RemoveSession(sessionId)
	this.DeleteCookie(c, w)
	
	// GAE の taskqueue.Lease() が正常に動作しないので治ってから実装する
//	_, err = taskqueue.LeaseByTag(c, 10, "Session", 100, "tag")
//	Check(c, err)
//	c.Debugf("Lease")
}

// Ajax のリクエストパラメータとして渡されたユーザを仮登録データベースに保存して、
// 本登録のためのメールを送信する。
// 24時間以内に登録されなかった場合のキャンセル処理を予約する。
func (this *Controller) InterimRegistration(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	model := NewModel(c)
	
	name := r.FormValue("name")
	mail := r.FormValue("mail")
	pass := r.FormValue("pass")
	
	if name == "" {
		fmt.Fprintf(w, `{"result":false, "msg":"ユーザ名が入力されていません"}`)
		return
	} else if mail == "" {
		fmt.Fprintf(w, `{"result":false, "msg":"メールアドレスが入力されていません"}`)
		return
	} else if pass == "" {
		fmt.Fprintf(w, `{"result":false, "msg":"パスワードが入力されていません"}`)
		return
	} else if model.ExistMail(mail) {
		fmt.Fprintf(w, `{"result":false, "msg":"既に登録されているメールアドレスです"}`)
		return
	}
	
	key := model.InterimRegistration(name, mail, pass)
	
	SendMail(c, "infomation@escape-3ds.appspotmail.com", mail, "仮登録完了のお知らせ", fmt.Sprintf(INTERIM_MAIL_BODY, name, key))
	
	// キャンセルタスクを予約
	values := url.Values{}
	values.Set("interim_key", key)
	delay, err := time.ParseDuration("24h")
	Check(c, err)

	task := taskqueue.NewPOSTTask("/cancel", values)
	task.Delay = delay
	taskqueue.Add(c, task, "default")
	
	fmt.Fprintf(w, `{"result":true}`)
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

// 仮登録のキャンセル
// 24時間以内に本登録されなかったらキャンセルする
// 既に登録されている場合は何もしない
func (this *Controller) Cancel(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	encodedInterimKey := r.FormValue("interim_key")
	
	model := NewModel(c)
	if(model.ExistInterimUser(encodedInterimKey)) {
		model.DeleteInterimUser(encodedInterimKey)
	}
	
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
	gameKey := model.AddGame(game, userKey)
	gameURL := game.ShortURL
	
	fmt.Fprintf(w, `{"key":"%s", "result":true, "name":"%s", "description":"%s", "url":"%s"}`, gameKey, gameName, gameDescription, gameURL)
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

// /sync/* にマッチしたら呼び出される。
// エディタで編集したゲーム情報をサーバと同期するための処理。
// URL を解析して更に細かいハンドラへ処理を割り振る。
func (this *Controller) SyncHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	path := strings.Split(r.URL.Path, "/") // "[0]/[1]sync/[2]kind/[3]id"
	
	model := NewModel(c)
	switch path[2] {
	case "game":
		model.SyncGame(w, r, path)
	case "scene":
		model.SyncScene(w, r, path)
	case "item":
		model.SyncItem(w, r, path)
	case "event":
		model.SyncEvent(w, r, path)
	}
}

// ゲームを JSON 形式に変換する
func (this *Controller) EncodeGame(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	encodedGameKey := r.FormValue("game_key")
	
	model := NewModel(c)
	encodedGame := model.GetGameJSON(encodedGameKey)
	
	fmt.Fprintf(w, `{"game":%s}`, encodedGame)
}

// イベントコードのアップデート
func (this *Controller) UpdateEventCode(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	code := r.FormValue("code")
	id := r.FormValue("id")
	
	model := NewModel(c)
	model.UpdateEventCode(id, code)
	fmt.Fprintf(w, `{}`)
}

// シーン開始時のイベントコードのアップデート
func (this *Controller) UpdateEnterCode(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	code := r.FormValue("code")
	sceneId := r.FormValue("id")
	
	model := NewModel(c)
	scene := model.GetScene(sceneId)
	scene.Enter = []byte(code)
	model.UpdateScene(sceneId, scene)
	
	fmt.Fprintf(w, `{}`)
}

// シーン終了時のイベントコードのアップデート
func (this *Controller) UpdateLeaveCode(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	code := r.FormValue("code")
	sceneId := r.FormValue("id")
	
	model := NewModel(c)
	scene := model.GetScene(sceneId)
	scene.Leave = []byte(code)
	model.UpdateScene(sceneId, scene)
	
	fmt.Fprintf(w, `{}`)
}

// ゲストアカウントの終了
func (this *Controller) ByeGuest(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	encodedUserKey := r.FormValue("user_key")
	
	model := NewModel(c)
	if model.ExistUser(encodedUserKey) {
		model.DeleteUser(encodedUserKey)
		this.Logout(w, r)
	}	
}

// パスワードの再発行
func (this *Controller) ResetPassword(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	mail := r.FormValue("mail")
	model := NewModel(c)
	key, user := model.GetUserFromMail(mail)
	rawPassword := user.ResetPassword(key, c)
	SendMail(c, "infomation@escape-3ds.appspotmail.com", user.Mail, "パスワードのリセットが完了しました", fmt.Sprintf(RESET_PASSWORD_MAIL_BODY, user.Name, rawPassword))
}

// パスワードの変更
func (this *Controller) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userKey := this.Session(w, r)

	c := appengine.NewContext(r)
	pass := r.FormValue("password")
	model := NewModel(c)
	model.ChangePassword(userKey, pass)
	fmt.Fprintf(w, `{}`)
}

// ゲーム名、ゲームの説明の変更
func (this *Controller) RenameGame(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	name := r.FormValue("name")
	description := r.FormValue("description")
	gameKey := r.FormValue("key")
	if name == "" {
		c.Warningf("ゲーム名を空文字列に変更しようとしました")
		return
	}
	model := NewModel(c)
	game := model.GetGame(gameKey)
	game.Name = name
	game.Description = description
	model.UpdateGame(gameKey, game)
	fmt.Fprintf(w, `{}`)
}

// 問合せ
func (this *Controller) Inquiry(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	body := r.FormValue("body")
	if body == "" {
		c.Warningf("メッセージ無しで問合せが実行されました")
		return
	}
	SendMail(c, "infomation@escape-3ds.appspotmail.com", ADMIN_MAIL, "ESCAPE-3DS にお問い合わせが来ました", fmt.Sprintf(body))
	fmt.Fprintf(w, `{}`)
}

// しばらく操作しなかった場合にセッションを終了する
// config.go の SESSION_TIME_LIMIT で時間を設定する
func (this *Controller) Timeout(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	sessionId := r.FormValue("session_id")
	this.CloseSession(c, sessionId, w, r)
}