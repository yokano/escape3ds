/**
 * データモデルの定義
 * @file
 */
package escape3ds

import (
	"appengine"
	"appengine/datastore"
	"appengine/memcache"
	"appengine/blobstore"
	"strings"
	"bytes"
	"fmt"
	"time"
	"encoding/json"
	"mime/multipart"
)

/**
 * モデル
 * @class
 * @property {appengine.Context} c コンテキスト
 */
type Model struct {
	c appengine.Context
}

/**
 * モデルの作成
 * @function
 * @param {appengine.Context} c コンテキスト
 * @returns {*Model} モデル
 */
func NewModel(c appengine.Context) *Model {
	model := new(Model)
	model.c = c
	return model
}

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
	if !exist([]string {"Twitter", "Facebook", "normal"}, data["user_type"]) {
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
	user.Pass, user.Salt = this.hashPassword(data["user_pass"], "")
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
 * ゲーム
 * @struct
 * @member {string} Name ゲーム名
 * @member {string} Description ゲームの説明
 * @member {string} Thumbnail サムネイルの画像パス
 * @member {string} UserKey 所有ユーザのエンコード済みキー
 * @member {string} FirstScene 最初のシーンのエンコード済みキー
 */
type Game struct {
	Name string
	Description string
	Thumbnail string
	UserKey string
	FirstScene string
}

/**
 * ゲームインスタンスの作成
 * @method
 * @memberof Model
 * @param {map[string]string}
 * {
 *     name: string
 *     description: string
 *     thumbnail: string
 *     user_key: string
 *     first_scene: string
 * }
 */
func (this *Model) NewGame(params map[string]string) *Game {
	game := new(Game)
	game.Name = params["name"]
	game.Description = params["description"]
	game.Thumbnail = params["thumbnail"]
	game.UserKey = params["user_key"]
	game.FirstScene = ""
	return game
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
func (this *Model) hashPassword(pass string, salt string) ([]byte, string) {
	if salt == "" {
		for i := 0; i < 4; i++ {
			salt = strings.Join([]string{salt, getRandomizedString()}, "")
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
func (this *Model) addUser(user *User) string {
	if user == nil {
		this.c.Errorf("ユーザの追加を中止しました")
		return ""
	}
	incompleteKey := datastore.NewIncompleteKey(this.c, "User", nil)
	completeKey, err := datastore.Put(this.c, incompleteKey, user)
	check(this.c, err)
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
func (this *Model) loginCheck(mail string, pass string) (string, string) {
	query := datastore.NewQuery("User").Filter("Mail =", mail)
	iterator := query.Run(this.c)
	key, err := iterator.Next(nil)
	
	if err != nil {
		this.c.Warningf("存在しないメールアドレスによるログインが試されました。アドレス：%s", mail)
		return "", ""
	}
	
	encodedKey := key.Encode()
	user := this.getUser(encodedKey)
	
	hashedPass, _ := this.hashPassword(pass, user.Salt)
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
func (this *Model) getUser(encodedKey string) *User {
	key, err := datastore.DecodeKey(encodedKey)
	check(this.c, err)
	
	user := new(User)
	err = datastore.Get(this.c, key, user)
	check(this.c, err)
	
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
func (this *Model) interimRegistration(name string, mail string, pass string) string {
	user := this.NewInterimUser(name, mail, pass)
	incompleteKey := datastore.NewIncompleteKey(this.c, "InterimUser", nil)
	completeKey, err := datastore.Put(this.c, incompleteKey, user)
	check(this.c, err)
	return completeKey.Encode()
}

/**
 * ユーザを本登録する
 * 仮登録データベース
 * @method
 * @memberof Model
 * @param {string} encodedKey エンコード済みの仮登録キー
 */
func (this *Model) registration(encodedKey string) {
	key, err := datastore.DecodeKey(encodedKey)
	check(this.c, err)
	
	interimUser := new(InterimUser)
	err = datastore.Get(this.c, key, interimUser)
	check(this.c, err)
	
	params := make(map[string]string, 5)
	params["user_type"] = "normal"
	params["user_name"] = interimUser.Name
	params["user_mail"] = interimUser.Mail
	params["user_pass"] = interimUser.Pass
	params["user_oauth_path"] = ""
	user := this.NewUser(params)
	this.addUser(user)
	
	err = datastore.Delete(this.c, key)
	check(this.c, err)
}

/**
 * 指定されたOAuthユーザが既にデータベース上に存在するかどうか調べる
 * @method
 * @memberof Model
 * @param {string} userType "Twitter"または"Facebook"
 * @param {string} oauthId 調べる対象のOAuthId
 * @returns {bool} 存在したらtrue
 */
func (this *Model) existOAuthUser(userType string, oauthId string) bool {
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
func (this *Model) getUserKey(params map[string]string) string {
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
 * データストアにゲームを追加する
 * @method
 * @memberof Model
 * @param {*Game} game 追加するゲーム
 * @returns {string} エンコード済みのゲームキー
 */
func (this *Model) addGame(game *Game) string {
	incompleteKey := datastore.NewIncompleteKey(this.c, "Game", nil)
	completeKey, err := datastore.Put(this.c, incompleteKey, game)
	check(this.c, err)
	return completeKey.Encode()
}

/**
 * データストアからゲームを取得する
 * @method
 * @memberof Model
 * @param {string} encodedGameKey エンコード済みのゲームキー
 * @returns {*Game} ゲームオブジェクト
 */
func (this *Model) getGame(encodedGameKey string) *Game {
	gameKey, err := datastore.DecodeKey(encodedGameKey)
	check(this.c, err)
	
	game := new(Game)
	err = datastore.Get(this.c, gameKey, game)
	check(this.c, err)
	
	return game
}

/**
 * データストアからゲームを削除する
 * 削除を命令したユーザとゲームの所有者が一致していることを事前に確認すること
 * この関数内ではチェックを行わない
 * @param {string} encodedGameKey エンコード済みのゲームキー
 */
func (this *Model) deleteGame(encodedGameKey string) {
	gameKey, err := datastore.DecodeKey(encodedGameKey)
	check(this.c, err)
	
	err = datastore.Delete(this.c, gameKey)
	check(this.c, err)
}

/**
 * ユーザが所有しているゲーム一覧を返す
 * @method
 * @memberof Model
 * @param {string} encodedUserKey ユーザキー
 * @returns {map[string]*Game} エンコード済みのゲームキーとゲームの対応表
 */
func (this *Model) getGameList(encodedUserKey string) map[string]*Game {
	query := datastore.NewQuery("Game").Filter("UserKey =", encodedUserKey)
	iterator := query.Run(this.c)
	
	count, err := query.Count(this.c)
	check(this.c, err)
	
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

/**
 * 仮登録ユーザ一覧を返す
 * @method
 * @memberof Model
 * @returns {map[string]*InterimUser} 仮登録ユーザリスト
 */
func (this *Model) getInterimUsers() map[string]*InterimUser {
	query := datastore.NewQuery("InterimUser")
	count, err := query.Count(this.c)
	check(this.c, err)
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
func (this *Model) getAllUser() map[string]*User {
	query := datastore.NewQuery("User")
	count, err := query.Count(this.c)
	check(this.c, err)
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

/**
 * セッションを開始する
 * memcache にセッションIDとユーザキーの対応を保存する
 * セッションは最後のページアクセスから24時間有効
 * 24時間経過したものは cron で定期的に削除される
 * @method
 * @memberof Model
 * @param {string}
 * @returns {string} セッションID
 */
func (this *Model) startSession(userKey string) string {
	sessionId := ""
	for i := 0; i < 4; i++ {
		sessionId = fmt.Sprintf("%s%s", sessionId, getRandomizedString())
	}
	expire := time.Now().Add(time.Hour * 24)
	
	data := make(map[string]string, 2)
	data["u"] = userKey
	data["e"] = expire.String()
	
	encodedData, err := json.Marshal(data)
	item := &memcache.Item {
		Key: sessionId,
		Value: encodedData,
	}
	err = memcache.Set(this.c, item)
	check(this.c, err)
	
	return sessionId
}

/**
 * memcache から指定されたセッション情報を削除する
 * @method
 * @memberof Model
 * @param {string} sessionId 対象のセッションID
 */
func (this *Model) removeSession(sessionId string) {
	err := memcache.Delete(this.c, sessionId)
	check(this.c, err)
}

/**
 * memcache からユーザキーを返す
 * 該当するユーザーキーが見つからなかったら空文字を返す
 * @method
 * @memberof Model
 * @param {string} sessionId セッションID
 * @returns {string} ユーザキー
 */
func (this *Model) getUserKeyFromSession(sessionId string) string {
	item, err := memcache.Get(this.c, sessionId)
	check(this.c, err)
	if err != nil {
		this.c.Warningf("セッションIDに関連付けられたユーザIDが存在しませんでした: sessionId:%s", sessionId)
		return ""
	}
	data := make(map[string]string, 2)
	err = json.Unmarshal(item.Value, &data)
	check(this.c, err)
	return data["u"]
}

/**
 * ファイルを blobstore に保存して blobkey を返す
 * @method
 * @memberof Model
 * @param {multipart.File} file 保存するファイル
 * @param {*multipart.FileHeader} fileHeader 保存するファイルのヘッダ
 */
func (this *Model) addBlob(file multipart.File, fileHeader *multipart.FileHeader) string {
	// ファイルデータを読み出す
	fileData := make([]byte, 256 * 1024)
	byteNum, err := file.Read(fileData)
	check(this.c, err)
	fileData = fileData[:byteNum]
	
	// blobstore に blob を作成
	mimeType := fileHeader.Header.Get("Content-Type")
	writer, err := blobstore.Create(this.c, mimeType)
	check(this.c, err)
	
	// blob にファイルデータを書き込む
	writedByteNum, err := writer.Write(fileData)
	check(this.c, err)
	if writedByteNum != byteNum {
		this.c.Errorf("blob の書き込みに失敗しました")
		return ""
	}
	err = writer.Close()
	check(this.c, err)
	
	// blob key を取得
	key, err := writer.Key()
	check(this.c, err)
	
	return string(key)
}

/**
 * blobstore に保存したファイルを読み出す
 * @method
 * @memberof Model
 * @param {string} key 読みだすファイルの blobkey
 * @returns {string} ファイルの Content-Type
 * @returns {[]byte} ファイルバイナリデータ
 */
func (this *Model) getBlob(key string) (string, []byte) {
	blobKey := appengine.BlobKey(key)
	
	// __BlobFileIndex__ から blobKey で Key Name (string id) を取得
	query := datastore.NewQuery("__BlobFileIndex__").Filter("blob_key =", blobKey)
	iterator := query.Run(this.c)
	blobFileIndexKey, err := iterator.Next(nil)
	check(this.c, err)
	
	// __Blobinfo__ から Key Name で size を取得
	query = datastore.NewQuery("__BlobInfo__").Filter("creation_handle =", blobFileIndexKey.StringID())
	iterator = query.Run(this.c)
	blobInfo := new(blobstore.BlobInfo)
	_, err = iterator.Next(blobInfo)
	check(this.c, err)
	size := blobInfo.Size
	
	// size バイトだけ reader からバイトを読み出す
	reader := blobstore.NewReader(this.c, blobKey)
	bytes := make([]byte, size)
	readed, err := reader.Read(bytes)
	check(this.c, err)
	if(int64(readed) != size) {
		this.c.Warningf("読み込んだファイルサイズとメタ情報のファイルサイズが異なります")
	}
	
	return blobInfo.ContentType, bytes
}