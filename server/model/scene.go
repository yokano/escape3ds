package model

import (
	"appengine/datastore"
	"net/http"
	"encoding/json"
	"fmt"
	. "server/lib"
)

// シーンオブジェクト。ゲームを構成する場面を表す。
// 背景画像、開始時のイベント、終了時のイベントを持つ。
type Scene struct {
	Name string
	Background string
	Enter string
	Leave string
}

// シーンオブジェクトを新しく作成する。
func (this *Model) NewScene(name string, background string) *Scene {
	scene := new(Scene)
	scene.Name = name
	scene.Background = background
	return scene
}

// シーンをデータストアへ追加する。
// シーンが追加されるのはエディタでシーンを新規作成した時だけ。
// Backbone.js がモデルに対して自動的に割り振る cid をデータストアの StringID として使う。
// 引数として追加するシーン、cid、追加先のゲームキーを渡す。
// 追加に成功したらエンコード済みのシーンキーを返す
func (this *Model) AddScene(scene *Scene, id string, encodedGameKey string) string {
	gameKey, err := datastore.DecodeKey(encodedGameKey)
	Check(this.c, err)
	sceneKey := datastore.NewKey(this.c, "Scene", id, 0, gameKey)
	_, err = datastore.Put(this.c, sceneKey, scene)
	encodedSceneKey := sceneKey.Encode()
	return encodedSceneKey
}

// 引数として指定された id のシーンをデータストアから削除する。
func (this *Model) DeleteScene(id string) {
	
}

// シーンデータの同期。
// /sync/scene/[gamekey] というURLでリクエストが送られる。
// リクエストのメソッドが CRUD に対応している。
// POST:CREATE, GET:READ, PUT:UPDATE, DELETE:DELETE
func (this *Model) SyncScene(w http.ResponseWriter, r *http.Request, path []string) {
	switch r.Method {
	case "POST":
		body := make([]byte, r.ContentLength)
		r.Body.Read(body)
		
		scene := new(Scene)
		json.Unmarshal(body, scene)
		
		gameKey, err := datastore.DecodeKey(path[3])
		Check(this.c, err)
		
		sceneKey := datastore.NewIncompleteKey(this.c, "Scene", gameKey)
		_, err = datastore.Put(this.c, sceneKey, scene)
		Check(this.c, err)
		
		encodedSceneKey := sceneKey.Encode()
		
		fmt.Fprintf(w, `{"sceneKey":"%s"}`, encodedSceneKey)
	case "GET":
	case "PUT":
	case "DELETE":
	}
}