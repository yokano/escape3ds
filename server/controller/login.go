package controller

import (
	"net/http"
	"appengine"
	"fmt"
	"net/url"
	"encoding/json"
	. "server/lib"
	. "server/model"
	. "server/view"
	. "server/config"
)

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
