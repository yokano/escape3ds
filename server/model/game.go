package model

import (
	"appengine/datastore"
	. "server/lib"
)

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
	Name string `json:"name"`
	Description string `json:"description"`
	Thumbnail string `json:"thumbnail"`
	UserKey string `json:"userKey"`
	FirstScene string `json:"firstScene"`
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
 * データストアにゲームを追加する
 * @method
 * @memberof Model
 * @param {*Game} game 追加するゲーム
 * @returns {string} エンコード済みのゲームキー
 */
func (this *Model) AddGame(game *Game) string {
	incompleteKey := datastore.NewIncompleteKey(this.c, "Game", nil)
	completeKey, err := datastore.Put(this.c, incompleteKey, game)
	Check(this.c, err)
	return completeKey.Encode()
}

/**
 * データストアからゲームを取得する
 * @method
 * @memberof Model
 * @param {string} encodedGameKey エンコード済みのゲームキー
 * @returns {*Game} ゲームオブジェクト
 */
func (this *Model) GetGame(encodedGameKey string) *Game {
	gameKey, err := datastore.DecodeKey(encodedGameKey)
	Check(this.c, err)
	
	game := new(Game)
	err = datastore.Get(this.c, gameKey, game)
	Check(this.c, err)
	
	return game
}

/**
 * ゲームを更新する
 * @method
 * @memberof Model
 * @param {string} encodedGameKey エンコード済みのゲームキー
 * @param {*Game} game 上書きするゲームオブジェクト
 */
func (this *Model) UpdateGame(encodedGameKey string, game *Game)  {
	gameKey, err := datastore.DecodeKey(encodedGameKey)
	Check(this.c, err)
	
	_, err = datastore.Put(this.c, gameKey, game)
	Check(this.c, err)
}

/**
 * データストアからゲームを削除する
 * 削除を命令したユーザとゲームの所有者が一致していることを事前に確認すること
 * この関数内ではチェックを行わない
 * @param {string} encodedGameKey エンコード済みのゲームキー
 */
func (this *Model) DeleteGame(encodedGameKey string) {
	gameKey, err := datastore.DecodeKey(encodedGameKey)
	Check(this.c, err)
	
	err = datastore.Delete(this.c, gameKey)
	Check(this.c, err)
}

/**
 * ユーザが所有しているゲーム一覧を返す
 * @method
 * @memberof Model
 * @param {string} encodedUserKey ユーザキー
 * @returns {map[string]*Game} エンコード済みのゲームキーとゲームの対応表
 */
func (this *Model) GetGameList(encodedUserKey string) map[string]*Game {
	query := datastore.NewQuery("Game").Filter("UserKey =", encodedUserKey)
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
