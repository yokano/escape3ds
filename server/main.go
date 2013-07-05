/**
 * エントリポイント
 * URLパターンから該当する処理へ振り分ける
 * 処理は controller.go に記載されている
 * @file
 */
package escape3ds

import "net/http"

/**
 * URLから処理を振り分ける
 * @function
 */
func init() {
	// 通常アクセス
	http.HandleFunc("/", top)
	http.HandleFunc("/editor", editor)
	http.HandleFunc("/gamelist", gamelist)
	http.HandleFunc("/logout", logout)
	
	// OAuth 関係
	http.HandleFunc("/login_twitter", loginTwitter)
	http.HandleFunc("/callback_twitter", callbackTwitter)
	http.HandleFunc("/login_facebook", loginFacebook)
	http.HandleFunc("/callback_facebook", callbackFacebook)

	// アカウント登録関係
	http.HandleFunc("/interim_registration", interimRegistration)
	http.HandleFunc("/registration", registration)

	// Ajax
	http.HandleFunc("/add_user", addUser)
	http.HandleFunc("/login", login)
	http.HandleFunc("/add_game", addGame)
	http.HandleFunc("/delete_game", deleteGame)
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/download", download)
	
	// 管理者専用 通常アクセス
	http.HandleFunc("/debug", debug)
	
	// 管理者専用 Ajax
	http.HandleFunc("/get_users", getUsers)
	http.HandleFunc("/get_interim_users", getInterimUsers)

	// エディタのセーブ
	http.HandleFunc("/sync/", syncHandler)
}