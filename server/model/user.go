package model

import (
	"strings"
	"appengine/datastore"
	"bytes"
	"fmt"
	. "server/lib"
)


// ユーザデータ
type User struct {
	Type string  // ユーザアカウントの種類 "Twitter"/"Facebook"/"normal"/"guest"
	Name string  // ユーザ名
	Pass []byte  // ユーザの暗号化済パスワード（user_type == "normal"の場合のみ）
	Mail string  // ユーザのメールアドレス（user_type == "normal"の場合のみ）
	Salt string  // パスワードソルト
	OAuthId string  // OAuthのサービスプロバイダが決めたユーザID
}

// ユーザオブジェクトを作成する。
// 戻り値は作成したユーザ、失敗したらnil。
func (this *Model) NewUser(data map[string]string) *User {
	// ユーザタイプチェック
	if !Exist([]string {"Twitter", "Facebook", "normal", "guest"}, data["user_type"]) {
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
	user.Pass, user.Salt = user.HashPassword(data["user_pass"], "")
	return user
}

// 仮登録ユーザ。トップページのユーザ登録フォームに入力されたユーザデータは一時的に仮登録データベースへ保存される。
// ユーザが登録したメールアドレスに借り登録完了のメールが送られ、そのなかのURLをクリックすることで、
// 仮登録データベースからユーザデータベースへデータが移行される。
// 仮登録完了から 24 時間以内にメールの URL がクリックされなかった場合は、
// 仮登録データベースからユーザ情報が削除され、仮登録は失敗と鳴る。
type InterimUser struct {
	Name string  // ユーザ名
	Mail string  // メールアドレス
	Pass string  // パスワード
}

// 仮登録ユーザを作成する。引数としてユーザ名、メールアドレス、パスワードを渡す。
// 新規作成された仮登録ユーザオブジェクトを返す。
func (this *Model) NewInterimUser(name string, mail string, pass string) *InterimUser {
	user := new(InterimUser)
	user.Name = name
	user.Mail = mail
	user.Pass = pass
	return user
}

// ユーザのパスワードをハッシュ化する。引数として平文のパスワードと、ソルトを渡す。
// 引数として渡すソルトが空文字の場合は自動でランダムな文字列を作成する。
// ハッシュ化が成功したら、ハッシュ化されたパスワードとソルトを返す。
// ハッシュは SHA-1 を使う。
func (this *User) HashPassword(pass string, salt string) ([]byte, string) {
	if salt == "" {
		for i := 0; i < 4; i++ {
			salt = strings.Join([]string{salt, GetRandomizedString()}, "")
		}
	}
	pass = strings.Join([]string{pass, salt}, "")
	hashedPass := SHA1(pass)
	return hashedPass, salt
}

// データストアに引数として渡したユーザを新規追加する。
// 追加が完了したらユーザキーを返す。
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

// ユーザの削除
func (this *Model) DeleteUser(encodedUserKey string) {
	if encodedUserKey == "" {
		this.c.Errorf("ユーザキーの指定なしで DeleteUser() が実行されました")
		return
	}
	userKey, err := datastore.DecodeKey(encodedUserKey)
	Check(this.c, err)
	err = datastore.Delete(this.c, userKey)
	Check(this.c, err)
}

// 引数として渡したメールアドレスとパスワードのユーザがいるか調べる。
// 存在した場合はユーザキーとユーザ名を返す。
// 存在しない場合は戻り値がすべて空文字になる。
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
	
	hashedPass, _ := user.HashPassword(pass, user.Salt)
	if bytes.Compare(user.Pass, hashedPass) != 0 {
		this.c.Warningf("間違ったパスワードが試されました。アドレス：%s", mail)
		return "", ""
	}
	
	return encodedKey, user.Name
}

// 引数として渡されたユーザキーのユーザをデータストアから取得する。
// 成功したら取得したユーザオブジェクトを返す。
func (this *Model) GetUser(encodedKey string) *User {
	key, err := datastore.DecodeKey(encodedKey)
	Check(this.c, err)
	
	user := new(User)
	err = datastore.Get(this.c, key, user)
	Check(this.c, err)
	
	return user
}


// ユーザを仮登録データベースに登録する。
// 仮登録されたユーザは24時間以内に本登録する。
// 本登録されなかった場合は24時間後に削除される。
// 仮登録ユーザのエンコードされたキーを返す。
func (this *Model) InterimRegistration(name string, mail string, pass string) string {
	user := this.NewInterimUser(name, mail, pass)
	incompleteKey := datastore.NewIncompleteKey(this.c, "InterimUser", nil)
	completeKey, err := datastore.Put(this.c, incompleteKey, user)
	Check(this.c, err)
	return completeKey.Encode()
}

// 仮登録ユーザキーを指定して、該当するユーザを登録する。
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

// 指定されたOAuthユーザが既にデータベース上に存在するかどうか調べる。
// 存在したら true、存在しなかったら false を返す。
// userType に "Twitter" または "Facebook" を指定すること。
// oauthId に調べる対象の OAuthId を指定する。
func (this *Model) ExistOAuthUser(userType string, oauthId string) bool {
	query := datastore.NewQuery("User").Filter("Type =", userType).Filter("OAuthId =", oauthId)
	iterator := query.Run(this.c)
	_, err := iterator.Next(nil)
	if err != nil {
		return false
	}
	return true
}

// パラメータで指定されたユーザを探してユーザキーを返す。
// 存在しない場合は空文字を返す。
// params には検索条件をセットする。
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

// 仮登録ユーザ一覧を返す
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

// 指定された仮登録ユーザがいるか調べる
func (this *Model) ExistInterimUser(encodedInterimKey string) bool {
	interimKey, err := datastore.DecodeKey(encodedInterimKey)
	Check(this.c, err)
	
	err = datastore.Get(this.c, interimKey, new(InterimUser))
	if(err == nil) {
		return true
	}
	Check(this.c, err)
	return false
}

// 仮登録ユーザを削除する
func (this *Model) DeleteInterimUser(encodedInterimKey string) {
	interimKey, err := datastore.DecodeKey(encodedInterimKey)
	Check(this.c, err)
	datastore.Delete(this.c, interimKey)
}

// すべてのユーザを取得する
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

// ユーザが所有しているゲーム一覧を返す。
// 戻り値は、エンコード済みのゲームキーとゲームの対応表
func (this *Model) GetGameList(encodedUserKey string) map[string]*Game {
	userKey, err := datastore.DecodeKey(encodedUserKey)
	Check(this.c, err)
	
	query := datastore.NewQuery("Game").Ancestor(userKey)
	iterator := query.Run(this.c)
	
	count, err := query.Count(this.c)
	Check(this.c, err)
	
	result := make(map[string]*Game, count)
	for ;; {
		game := new(Game)
		gameKey, err := iterator.Next(game)
		if err != nil {
			break
		}
		result[gameKey.Encode()] = game
	}
	
	return result
}
