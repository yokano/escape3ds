package model

import (
	"appengine/datastore"
	"net/http"
	"encoding/json"
	"fmt"
	. "server/lib"
	. "server/config"
)

// ゲームオブジェクト。ゲームはそれを所有するユーザにぶら下がる形で定義される。
// ゲームは自分を所有するユーザを親として参照する形で datastore に保存される。
// datastore の ancestor path を参照。
// https://developers.google.com/appengine/docs/go/datastore/#Ancestor_Paths
type Game struct {
	Name string `json:"name"`  // ゲーム名
	Description string `json:"description"`  // ゲームの説明
	Thumbnail string `json:"thumbnail"`  // サムネイルの画像パス
	FirstScene string `json:"firstScene"`  // 最初のシーンのエンコード済みキー
	SceneList map[string]*Scene `json:"sceneList" datastore:"-"`  // シーン一覧
	ItemList map[string]*Item `json:"itemList" datastore:"-"`  // アイテム一覧
	ShortURL string `json:"shortURL"`  // 短縮URL
}

// 新しいゲームオブジェクトを作成して返す。新しく作られるゲームのパラメータを格納した map 引数として渡す。
func (this *Model) NewGame(params map[string]string) *Game {
	game := new(Game)
	game.Name = params["name"]
	game.Description = params["description"]
	game.Thumbnail = params["thumbnail"]
	game.FirstScene = ""
	return game
}

// データストアにゲームを追加してエンコード済みのキーを返す。
// 引数として追加するゲームオブジェクトと、エンコード済みユーザキーを渡す。
// 引数として渡したユーザキーにぶら下がる形でゲームが保存される。
func (this *Model) AddGame(game *Game, encodedUserKey string) string {
	userKey, err := datastore.DecodeKey(encodedUserKey)
	Check(this.c, err)
	
	incompleteKey := datastore.NewIncompleteKey(this.c, "Game", userKey)
	completeKey, err := datastore.Put(this.c, incompleteKey, game)
	Check(this.c, err)
	
	// 短縮URLを取得する
	key := GOOGLE_API_KEY
	longUrl := Join("http://", HOSTNAME, "/runtime?game_key=", completeKey.Encode())
	requestBody := Join(`{"key":"`, key, `","longUrl":"`, longUrl, `"}`)
	params := make(map[string]string, 1)
	params["Content-Type"] = "application/json"
	response := Request(this.c, "POST", "https://www.googleapis.com/urlshortener/v1/url", params, requestBody)

	responseBody := make([]byte, response.ContentLength)
	_, err = response.Body.Read(responseBody)
	Check(this.c, err)

	result := make(map[string]string, 3)
	err = json.Unmarshal(responseBody, &result)
	Check(this.c, err)
	
	game.ShortURL = result["id"]
	this.UpdateGame(completeKey.Encode(), game)
	
	return completeKey.Encode()
}

// エンコード済みのキーを指定してデータストアからゲームを取得する。
func (this *Model) GetGame(encodedGameKey string) *Game {
	gameKey, err := datastore.DecodeKey(encodedGameKey)
	Check(this.c, err)
	
	game := new(Game)
	err = datastore.Get(this.c, gameKey, game)
	Check(this.c, err)
	
	game.SceneList = this.GetScenes(encodedGameKey)
	game.ItemList = this.GetItems(encodedGameKey)
	
	return game
}

// JSON 形式でゲームを取得する
func (this *Model) GetGameJSON(encodedGameKey string) string {
	game := this.GetGame(encodedGameKey)
	game.ItemList = this.GetItems(encodedGameKey)
	game.SceneList = this.GetScenes(encodedGameKey)

	for key, val := range game.SceneList {
		val.EventList = this.GetEventList(key)
	}
	
	bytes, err := json.Marshal(game)
	Check(this.c, err)
	
	return string(bytes)
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
	game := this.GetGame(encodedGameKey)
	
	// シーン削除
	for sceneKey, _ := range game.SceneList {
		this.DeleteScene(sceneKey)
	}
	
	// アイテム削除
	for itemKey, _ := range game.ItemList {
		this.DeleteItem(itemKey)
	}
	
	// ゲーム削除
	gameKey, err := datastore.DecodeKey(encodedGameKey)
	Check(this.c, err)
	err = datastore.Delete(this.c, gameKey)
	Check(this.c, err)
}

// ゲーム内のシーン一覧を取得する。引数としてゲームキーを渡す。
// 戻り値としてシーンID(エンコード済みキー)を key、シーンオブジェクトを value とした map を返す。
func (this *Model) GetScenes(encodedGameKey string) map[string]*Scene {
	gameKey := DecodeKey(this.c, encodedGameKey)
	
	query := datastore.NewQuery("Scene").Ancestor(gameKey)
	count, err := query.Count(this.c)
	Check(this.c, err)
	
	iterator := query.Run(this.c)
	scenes := make(map[string]*Scene, count)
	
	for ;; {
		sceneKey, err := iterator.Next(nil)
		if err != nil {
			break
		}
		encodedSceneKey := sceneKey.Encode()
		scene := this.GetScene(encodedSceneKey)
		scenes[encodedSceneKey] = scene
	}
	
	return scenes
}

// 指定されたシーンが含まれているゲームを取得する
func (this *Model) GetGameFromScene(encodedSceneKey string) (string, *Game) {
	sceneKey, err := datastore.DecodeKey(encodedSceneKey)
	Check(this.c, err)
	
	gameKey := sceneKey.Parent()
	encodedGameKey := gameKey.Encode()
	
	return encodedGameKey, this.GetGame(encodedGameKey)
}

// ゲームデータの同期。
// リクエストのメソッドが CRUD に対応している。
// POST:CREATE, GET:READ, PUT:UPDATE, DELETE:DELETE
func (this *Model) SyncGame(w http.ResponseWriter, r *http.Request, path []string) {
	switch r.Method {
	case "POST":
	case "GET":
	case "PUT":
		gameKey := path[3]
		body := make([]byte, r.ContentLength)
		r.Body.Read(body)
		
		game := new(Game)
		json.Unmarshal(body, game)
		
		model := NewModel(this.c)
//		oldGame := model.GetGame(gameKey)
//		if oldGame.UserKey != this.Session(w, r) {
//			c.Warningf("他者のゲームを削除しようとしました　userKey:%s, gameKey:%s", this.Session(w, r), gameKey)
//			return
//		}
		
		model.UpdateGame(gameKey, game)
		fmt.Fprintf(w, `{}`)
	case "DELETE":
		this.c.Debugf("DELETE GAME")
	}
}
