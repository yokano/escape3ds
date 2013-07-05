/**
 * HTMLの表示
 * @file
 */
package escape3ds

import(
	"net/http"
	"html/template"
	"appengine"
	. "server/model"
	. "server/lib"
)

/**
 * 画面表示を行うクラス
 * @class
 * @property {appengine.Context} c コンテキスト
 * @property {http.ResponseWriter} w 応答先
 * @property {*http.Request} r リクエスト
 */
type View struct {
	c appengine.Context
	w http.ResponseWriter
}

/**
 * View の作成
 * @function
 * @param {appengine.Context} c コンテキスト
 * @param {http.ResponseWriter} w 応答先
 * @returns {*View} 作成したView
 */
func NewView(c appengine.Context, w http.ResponseWriter) *View {
	view := new(View)
	view.c = c
	view.w = w
	return view
}

/**
 * ログイン画面を表示する
 * @method
 * @memberof View
 */
func (this *View) login() {
	t, err := template.ParseFiles("server/html/login.html")
	Check(this.c, err)
	t.Execute(this.w, nil)
}

/**
 * エディタ画面を表示する
 * @method
 * @memberof View
 * @param {string} gameKey ゲームキー
 */
func (this *View) editor(gameKey string) {
	model := NewModel(this.c)
	game := model.GetGame(gameKey)
	
	type Contents struct {
		Game *Game
		GameKey string
	}
	contents := new(Contents)
	contents.Game = game
	contents.GameKey = gameKey
	
	t, err := template.ParseFiles("server/html/editor.html")
	Check(this.c, err)
	t.Execute(this.w, contents)
}

/**
 * デバッグ画面の表示
 * @method
 * @memberof View
 */
func (this *View) debug() {
	t, err := template.ParseFiles("server/html/debug.html")
	Check(this.c, err)
	t.Execute(this.w, nil)
}

/**
 * 仮登録ページの表示
 * @method
 * @memberof View
 */
func (this *View) interimRegistration() {
	t, err := template.ParseFiles("server/html/interim_registration.html")
	Check(this.c, err)
	t.Execute(this.w, nil)
}

/**
 * 本登録完了ページの表示
 * @method
 * @memberof View
 */
func (this *View) registration() {
	t, err := template.ParseFiles("server/html/registration.html")
	Check(this.c, err)
	t.Execute(this.w, nil)
}

/**
 * ゲーム一覧の表示
 * @method
 * @memberof View
 * @param {string} userKey
 */
func (this *View) gamelist(userKey string) {
	model := NewModel(this.c)
	gameList := model.GetGameList(userKey)
	
	t, err := template.ParseFiles("server/html/gamelist.html")
	Check(this.c, err)
	t.Execute(this.w, gameList)
}