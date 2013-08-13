// 画面表示全般を行う。必要に応じてモデルからデータを取り出し整形する。
// HTML ファイルを出力するときは html ディレクトリからテンプレートを取り出して使う。
package view

import(
	"net/http"
	"html/template"
	"appengine"
	. "server/model"
	. "server/lib"
)

// 画面表示を行うオブジェクト
type View struct {
	c appengine.Context
	w http.ResponseWriter
}

// View オブジェクトを作成する。
func NewView(c appengine.Context, w http.ResponseWriter) *View {
	view := new(View)
	view.c = c
	view.w = w
	return view
}

// ログイン画面を表示する
func (this *View) Login() {
	t, err := template.ParseFiles("server/view/html/login.html")
	Check(this.c, err)
	t.Execute(this.w, nil)
}

// エディタ画面を表示する。エディットするゲームのキーを引数として渡す。
func (this *View) Editor(encodedGameKey string) {
	model := NewModel(this.c)
	game := model.GetGame(encodedGameKey)
	
	// ゲームやシーンの初期状態を JavaScript へ埋め込む
	type Contents struct {
		Game *Game
		GameKey string
		Scenes map[string]*Scene
		Items map[string]*Item
	}
	contents := new(Contents)
	contents.Game = game
	contents.GameKey = encodedGameKey
	contents.Scenes = model.GetScenes(encodedGameKey)
	contents.Items = model.GetItems(encodedGameKey)
	
	t, err := template.ParseFiles("server/view/html/editor.html")
	Check(this.c, err)
	
	t.Execute(this.w, contents)
}

// デバッグ画面の表示。デバッグ画面ではユーザの登録をメールを介さずに行ったりするための画面。
func (this *View) Debug() {
	t, err := template.ParseFiles("server/view/html/debug.html")
	Check(this.c, err)
	t.Execute(this.w, nil)
}

// 仮登録ページの表示。ユーザがトップページでメールアドレスとパスワードを入力した後に表示される。
// 本登録のためのメールアドレスが送信される。
func (this *View) InterimRegistration() {
	t, err := template.ParseFiles("server/view/html/interim_registration.html")
	Check(this.c, err)
	t.Execute(this.w, nil)
}

// 本登録完了ページの表示。仮登録状態のユーザの元へ送られたメールからジャンプしてくる。
// このページが表示されたら登録が完了。
func (this *View) Registration() {
	t, err := template.ParseFiles("server/view/html/registration.html")
	Check(this.c, err)
	t.Execute(this.w, nil)
}

// ゲーム一覧の表示。ログイン後に表示される最初のページ。
// 引数としてユーザキーを渡す。
// セッションIDがブラウザに保存されている場合は、トップページではなくこちらが表示される。
func (this *View) Gamelist(userKey string) {
	model := NewModel(this.c)
	gameList := model.GetGameList(userKey)
	
	t, err := template.ParseFiles("server/view/html/gamelist.html")
	Check(this.c, err)
	t.Execute(this.w, gameList)
}

// テストプレイ画面の表示
func (this *View) Runtime(gameKey string) {
	model := NewModel(this.c)
	game := model.GetGameJSON(gameKey)
	
	t, err := template.ParseFiles("server/view/html/runtime.html")
	Check(this.c, err)
	t.Execute(this.w, game)
}

// イベントエディタの表示（クリック時）
func (this *View) EventEditor(gameKey string, eventKey string) {
	model := NewModel(this.c)
	event := model.GetEvent(eventKey)
	itemList := model.GetItems(gameKey)
	sceneList := model.GetScenes(gameKey)
	
	contents := make(map[string]interface{}, 0)
	contents["trigger"] = "click"
	contents["id"] = eventKey
	contents["code"] = event.RawCode
	contents["itemList"] = itemList
	contents["sceneList"] = sceneList
	
	t, err := template.ParseFiles("server/view/html/event_editor.html")
	Check(this.c, err)
	t.Execute(this.w, contents)
}

// イベントエディタの表示（シーン開始時）
func (this *View) EnterEventEditor(gameKey string, sceneKey string) {
	model := NewModel(this.c)
	scene := model.GetScene(sceneKey)
	itemList := model.GetItems(gameKey)
	sceneList := model.GetScenes(gameKey)
	
	contents := make(map[string]interface{}, 0)
	contents["trigger"] = "enter"
	contents["id"] = sceneKey
	contents["code"] = scene.RawEnter
	contents["itemList"] = itemList
	contents["sceneList"] = sceneList
	
	t, err := template.ParseFiles("server/view/html/event_editor.html")
	Check(this.c, err)
	t.Execute(this.w, contents)
}