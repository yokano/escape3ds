// OAuth 2.0 による通信。
// authorization code 方式のみ。
package lib

import (
	"appengine"
	"net/http"
	"strings"
	"fmt"
)

// OAuth 2.0
type OAuth2 struct {
	context appengine.Context
	clientId string
	clientSecret string
}


// OAuth2.0 インスタンスを返す
func NewOAuth2(c appengine.Context, clientId string, clientSecret string) *OAuth2 {
	oauth := new(OAuth2)
	oauth.context = c
	oauth.clientId = clientId
	
	oauth.clientSecret = clientSecret
	return oauth
}

// 認証コードをリクエストする
// 認証ページヘのリダイレクトを行いユーザに認証してもらう
// 認証が完了したらリダイレクトURIへ認証コードが返ってくる
func (this *OAuth2) RequestAuthorizationCode(w http.ResponseWriter, r *http.Request, targetUri string, redirectUri string) {
	targetUri = fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code", targetUri, this.clientId, redirectUri)
	http.Redirect(w, r, targetUri, 302)
}

// アクセストークンをリクエストする
// 引き換えとして認証コードを渡すこと
func (this *OAuth2) RequestAccessToken(w http.ResponseWriter, r *http.Request, targetUri string, redirectUri string, code string) string {
	targetUri = fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&client_secret=%s&code=%s", targetUri, this.clientId, redirectUri, this.clientSecret, code)
	
	params := make(map[string]string, 4)
	params["client_id"] = this.clientId
	params["redirect_uri"] = redirectUri
	params["client_secret"] = this.clientSecret
	params["code"] = code
	response := Request(this.context, "GET", targetUri, params, "")
	
	body := make([]byte, 1024)
	_, err := response.Body.Read(body)
	Check(this.context, err)
	
	// response: oauth_token=******&expires=******
	responseParams := strings.Split(string(body), "&")
	tokenParam := strings.Split(responseParams[0], "=")
	return tokenParam[1]
}

// アクセストークンを使ってAPIを呼び出す
func (this *OAuth2) RequestAPI(w http.ResponseWriter, targetUri string, accessToken string) []byte {
	params := make(map[string]string, 1)
	params["access_token"] = accessToken
	response := Request(this.context, "GET", targetUri, params, "")
	
	result := make([]byte, 1024)
	i, err := response.Body.Read(result)
	Check(this.context, err)
	result = result[0:i]
	return result
}