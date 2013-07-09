package model

import (
	"appengine/datastore"
	. "server/lib"
)

// ゲームオブジェクト。ゲームはそれを所有するユーザにぶら下がる形で定義される。
// ゲームは自分を所有するユーザを親として参照する形で datastore に保存される。
// datastore の ancestor path を参照。
// https://developers.google.com/appengine/docs/go/datastore/#Ancestor_Paths
type Game struct {
	Name string `json:"name"`  // ゲーム名
	Description string `json:"description"`  // ゲームの説明
	Thumbnail string `json:"thumbnail"`  // サムネイルの画像パス
	UserKey string `json:"userKey"`  // 所有ユーザのエンコード済みキー
	FirstScene string `json:"firstScene"`  // 最初のシーンのエンコード済みキー
}

// 新しいゲームオブジェクトを作成して返す。新しく作られるゲームのパラメータを格納した map 引数として渡す。
func (this *Model) NewGame(params map[string]string) *Game {
	game := new(Game)
	game.Name = params["name"]
	game.Description = params["description"]
	game.Thumbnail = params["thumbnail"]
	game.UserKey = params["user_key"]
	game.FirstScene = ""
	return game
}

// データストアにゲームを追加してエンコード済みのキーを返す。
func (this *Model) AddGame(game *Game) string {
	incompleteKey := datastore.NewIncompleteKey(this.c, "Game", nil)
	completeKey, err := datastore.Put(this.c, incompleteKey, game)
	Check(this.c, err)
	return completeKey.Encode()
}

// エンコード済みのキーを指定してデータストアからゲームを取得する。
func (this *Model) GetGame(encodedGameKey string) *Game {
	gameKey, err := datastore.DecodeKey(encodedGameKey)
	Check(this.c, err)
	
	game := new(Game)
	err = datastore.Get(this.c, gameKey, game)
	Check(this.c, err)
	
	return game
}

// ゲームを更新する。引数として渡されたキー encodedGameKey のゲームを、ゲームオブジェクト game で上書きする。
func (this *Model) UpdateGame(encodedGameKey string, game *Game)  {
	gameKey, err := datastore.DecodeKey(encodedGameKey)
	Check(this.c, err)
	
	_, err = datastore.Put(this.c, gameKey, game)
	Check(this.c, err)
}

// データストアから指定したキーのゲームを削除する。
// 削除を命令したユーザとゲームの所有者が一致していることを事前に確認すること。
// この関数内ではチェックを行わない。
func (this *Model) DeleteGame(encodedGameKey string) {
	gameKey, err := datastore.DecodeKey(encodedGameKey)
	Check(this.c, err)
	
	err = datastore.Delete(this.c, gameKey)
	Check(this.c, err)
}

// ユーザが所有しているゲーム一覧を返す。
// 戻り値は、エンコード済みのゲームキーとゲームの対応表
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
