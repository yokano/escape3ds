package escape3ds

import (
	"strings"
	"appengine/datastore"
	"bytes"
	"fmt"
	. "server/lib"
)

/**
 * ユーザデータ
 * @struct
 * @property {string} Type ユーザアカウントの種類 "Twitter"/"Facebook"/"normal"
 * @property {string} Name ユーザ名
 * @property {[]byte} Pass ユーザの暗号化済パスワード（user_type == "normal"の場合のみ）
 * @property {string} Mail ユーザのメールアドレス（user_type == "normal"の場合のみ）
 * @property {string} OAuthId OAuthのサービスプロバイダが決めたユーザID
 */
type User struct {
	Type string
	Name string
	Pass []byte
	Mail string
	Salt string
	OAuthId string
}

/**
 * ユーザを作成する
 * @method
 * @memberof Model
 * @param {map[string]string} ユーザの設定項目を含んだマップ
 * {
 *     user_type: string
 *     user_name: string
 *     user_mail: string
 *     user_oauth_id: string
 *     user_pass: string
 * }
 * @returns {*User} ユーザ、失敗したらnil
 */
func (this *Model) NewUser(data map[string]string) *User {
	// ユーザタイプチェック
	if !Exist([]string {"Twitter", "Facebook", "normal"}, data["user_type"]) {
		this.c.Errorf("不正なユーザタイプが入力されました")
		return nil
	}
	
	// OAuthアカウントチェック
	if data["user_type"] == "Twitter" || data["user_type"] == "Facebook" {
		if data["user_oauth_id"] == "" {
			this.c.Errorf("OAuthアカウントのidが設定されていません")
			return nil
		}
	}
	
	// 通常アカウントチェック
	if data["user_type"] == "normal"{
		if data["user_mail"] == "" {
			this.c.Errorf("メールアドレスが入力されていません")
			return nil
		}
		if data["user_pass"] == "" {
			this.c.Errorf("パスワードが入力されていません")
			return nil
		}
	}
	
	user := new(User)
	user.Type = data["user_type"]
	user.Name = data["user_name"]
	user.Mail = data["user_mail"]
	user.OAuthId = data["user_oauth_id"]
	user.Pass, user.Salt = this.HashPassword(data["user_pass"], "")
	return user
}

/**
 * 仮登録ユーザ
 * @struct
 * @member {string} Name ユーザ名
 * @member {string} Mail メールアドレス
 * @member {string} Pass パスワード
 */
type InterimUser struct {
	Name string
	Mail string
	Pass string
}

/**
 * 仮登録ユーザの作成
 * @function
 * @param {string} name ユーザ名
 * @param {string} mail メールアドレス
 * @param {string} pass パスワード
 * @returns {*InterimUser} 仮登録ユーザ
 */
func (this *Model) NewInterimUser(name string, mail string, pass string) *InterimUser {
	user := new(InterimUser)
	user.Name = name
	user.Mail = mail
	user.Pass = pass
	return user
}

/**
 * ユーザのパスワードをハッシュ化する
 * @method
 * @memberof Model
 * @param {string} pass 平文パスワード
 * @param {string} salt ソルト。空文字が渡された場合は自動で作成する。
 * @returns {[]byte} 暗号化されたパスワード
 * @returns {string} 使用したソルト
 */
func (this *Model) HashPassword(pass string, salt string) ([]byte, string) {
	if salt == "" {
		for i := 0; i < 4; i++ {
			salt = strings.Join([]string{salt, GetRandomizedString()}, "")
		}
	}
	pass = strings.Join([]string{pass, salt}, "")
	hashedPass := SHA1(pass)
	return hashedPass, salt
}

/**
 * ユーザの追加
 * @method
 * @memberof Model
 * @param {*User} user 追加するユーザ
 * @returns {string} エンコードされたユーザキー
 */
func (this *Model) AddUser(user *User) string {
	if user == nil {
		this.c.Errorf("ユーザの追加を中止しました")
		return ""
	}
	incompleteKey := datastore.NewIncompleteKey(this.c, "User", nil)
	completeKey, err := datastore.Put(this.c, incompleteKey, user)
	Check(this.c, err)
	encodedKey := completeKey.Encode()
	return encodedKey
}

/**
 * 指定されたメールアドレスとパスワードのユーザがいるか調べる
 * 存在しない場合は戻り値がすべて空文字になる
 * @method
 * @memberof Member
 * @returns {string} エンコードされたキー
 * @returns {string} ユーザ名
 */
func (this *Model) LoginCheck(mail string, pass string) (string, string) {
	query := datastore.NewQuery("User").Filter("Mail =", mail)
	iterator := query.Run(this.c)
	key, err := iterator.Next(nil)
	
	if err != nil {
		this.c.Warningf("存在しないメールアドレスによるログインが試されました。アドレス：%s", mail)
		return "", ""
	}
	
	encodedKey := key.Encode()
	user := this.GetUser(encodedKey)
	
	hashedPass, _ := this.HashPassword(pass, user.Salt)
	if bytes.Compare(user.Pass, hashedPass) != 0 {
		this.c.Warningf("間違ったパスワードが試されました。アドレス：%s", mail)
		return "", ""
	}
	
	return encodedKey, user.Name
}

/**
 * ユーザの取得
 * @method
 * @memberof Model
 * @param {string} encodedKey エンコードされたキー
 */
func (this *Model) GetUser(encodedKey string) *User {
	key, err := datastore.DecodeKey(encodedKey)
	Check(this.c, err)
	
	user := new(User)
	err = datastore.Get(this.c, key, user)
	Check(this.c, err)
	
	return user
}

/**
 * ユーザを仮登録する
 * 仮登録したユーザは24時間以内に本登録する
 * 本登録されなかった場合は24時間後に削除される
 * @method
 * @memberof Model
 * @param {string} name ユーザ名
 * @param {string} mail メールアドレス
 * @param {string} pass パスワード
 * @returns {string} 仮登録ユーザのエンコードされたキー
 */
func (this *Model) InterimRegistration(name string, mail string, pass string) string {
	user := this.NewInterimUser(name, mail, pass)
	incompleteKey := datastore.NewIncompleteKey(this.c, "InterimUser", nil)
	completeKey, err := datastore.Put(this.c, incompleteKey, user)
	Check(this.c, err)
	return completeKey.Encode()
}

/**
 * ユーザを本登録する
 * 仮登録データベース
 * @method
 * @memberof Model
 * @param {string} encodedKey エンコード済みの仮登録キー
 */
func (this *Model) Registration(encodedKey string) {
	key, err := datastore.DecodeKey(encodedKey)
	Check(this.c, err)
	
	interimUser := new(InterimUser)
	err = datastore.Get(this.c, key, interimUser)
	Check(this.c, err)
	
	params := make(map[string]string, 5)
	params["user_type"] = "normal"
	params["user_name"] = interimUser.Name
	params["user_mail"] = interimUser.Mail
	params["user_pass"] = interimUser.Pass
	params["user_oauth_path"] = ""
	user := this.NewUser(params)
	this.AddUser(user)
	
	err = datastore.Delete(this.c, key)
	Check(this.c, err)
}

/**
 * 指定されたOAuthユーザが既にデータベース上に存在するかどうか調べる
 * @method
 * @memberof Model
 * @param {string} userType "Twitter"または"Facebook"
 * @param {string} oauthId 調べる対象のOAuthId
 * @returns {bool} 存在したらtrue
 */
func (this *Model) ExistOAuthUser(userType string, oauthId string) bool {
	query := datastore.NewQuery("User").Filter("Type =", userType).Filter("OAuthId =", oauthId)
	iterator := query.Run(this.c)
	_, err := iterator.Next(nil)
	if err != nil {
		return false
	}
	return true
}

/**
 * パラメータで指定されたユーザを探してキーを返す
 * 存在しない場合は空文字を返す
 * @method
 * @memberof Model
 * @param {map[string]string} 検索条件
 * @returns {string} 該当するユーザのキー、または空文字
 */
func (this *Model) GetUserKey(params map[string]string) string {
	query := datastore.NewQuery("User")
	for pkey, pval := range params {
		query.Filter(fmt.Sprintf("%s =", pkey), pval)
	}
	iterator := query.Run(this.c)
	key, err := iterator.Next(nil)
	if err != nil {
		return ""
	}
	return key.Encode()
}

/**
 * 仮登録ユーザ一覧を返す
 * @method
 * @memberof Model
 * @returns {map[string]*InterimUser} 仮登録ユーザリスト
 */
func (this *Model) GetInterimUsers() map[string]*InterimUser {
	query := datastore.NewQuery("InterimUser")
	count, err := query.Count(this.c)
	Check(this.c, err)
	iterator := query.Run(this.c)
	result := make(map[string]*InterimUser, count)
	for i := 0; i < count; i++ {
		user := new(InterimUser)
		key, err := iterator.Next(user)
		if err != nil {
			break
		}
		encodedKey := key.Encode()
		result[encodedKey] = user
	}
	return result
}

/**
 * すべてのユーザを取得する
 * @method
 * @memberof Model
 * @returns {map[string]*User} ユーザ一覧
 */
func (this *Model) GetAllUser() map[string]*User {
	query := datastore.NewQuery("User")
	count, err := query.Count(this.c)
	Check(this.c, err)
	iterator := query.Run(this.c)
	result := make(map[string]*User, count)
	for i := 0; i < count; i++ {
		user := new(User)
		key, err := iterator.Next(user)
		if err != nil {
			break
		}
		encodedKey := key.Encode()
		result[encodedKey] = user
	}
	return result
}
