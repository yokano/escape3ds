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
	Name string        // シーン名
	Background string  // 背景画像のBlobkey、設定していなければ空文字
	Enter string       // シーン開始時のイベント
	Leave string       // シーン終了時のイベント
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

// id で指定されたシーンを新しいシーン scene で更新する。
func (this *Model) UpdateScene(id string, scene *Scene) {
	
}

// ゲーム内のシーン一覧を取得する。引数としてゲームキーを渡す。
// 戻り値としてシーンID(エンコード済みキー)を key、シーンオブジェクトを value とした map を返す。
func (this *Model) GetScenes(encodedGameKey string) map[string]*Scene {
	gameKey, err := datastore.DecodeKey(encodedGameKey)
	Check(this.c, err)
	
	query := datastore.NewQuery("Scene").Ancestor(gameKey)
	count, err := query.Count(this.c)
	Check(this.c, err)
	
	iterator := query.Run(this.c)
	scenes := make(map[string]*Scene, count)
	
	for ;; {
		scene := new(Scene)
		sceneKey, err := iterator.Next(scene)
		if err != nil {
			break
		}
		encodedSceneKey := sceneKey.Encode()
		scenes[encodedSceneKey] = scene
	}
	
	return scenes
}

// シーンデータの同期処理を行う。
// /sync/scene/[gamekey] というURLでリクエストが送られる。
// リクエストのメソッドが CRUD に対応している。
// POST:CREATE, GET:READ, PUT:UPDATE, DELETE:DELETE.
// リクエストのメソッドに合わせて適切な関数を使用する。
func (this *Model) SyncScene(w http.ResponseWriter, r *http.Request, path []string) {
	switch r.Method {
	case "POST":
		body := GetRequestBodyJSON(r)
		scene := new(Scene)
		json.Unmarshal(body, scene)
		
		gameKey, err := datastore.DecodeKey(path[3])
		Check(this.c, err)
		
		sceneKey := datastore.NewIncompleteKey(this.c, "Scene", gameKey)
		sceneKey, err = datastore.Put(this.c, sceneKey, scene)
		Check(this.c, err)
		
		encodedSceneKey := sceneKey.Encode()
		
		fmt.Fprintf(w, `{"sceneKey":"%s"}`, encodedSceneKey)

	case "PUT":
		body := GetRequestBodyJSON(r)
		newScene := new(Scene)
		json.Unmarshal(body, newScene)
		
		encodedSceneKey := path[4]
		sceneKey, err := datastore.DecodeKey(encodedSceneKey)
		Check(this.c, err)
		
		oldScene := new(Scene)
		err = datastore.Get(this.c, sceneKey, oldScene)
		Check(this.c, err)
		
		if oldScene.Background != "" && oldScene.Background != newScene.Background && !this.CheckDuplicateBackground(oldScene.Background) {
			this.DeleteBlob(oldScene.Background)
		}
		
		sceneKey, err = datastore.Put(this.c, sceneKey, newScene)
		Check(this.c, err)
		
		fmt.Fprintf(w, `{}`)

	case "DELETE":
		encodedSceneKey := path[4]
		sceneKey, err := datastore.DecodeKey(encodedSceneKey)
		Check(this.c, err)
		
		scene := new(Scene)
		err = datastore.Get(this.c, sceneKey, scene)
		Check(this.c, err)
		
		// 他のシーンから画像が参照されていなければ画像を削除する
		if scene.Background != "" && !this.CheckDuplicateBackground(scene.Background) {
			this.DeleteBlob(scene.Background)
		}
		
		err = datastore.Delete(this.c, sceneKey)
		Check(this.c, err)

		fmt.Fprintf(w, `{}`)
	}
}

// 指定された背景画像が２つ以上のシーンから参照されているかどうかを調べる。
// 複数から参照されていれば ture, それ以外は false を返す。
// background は blobkey。
func (this *Model) CheckDuplicateBackground(background string) bool {
	q := datastore.NewQuery("Scene").Filter("Background =", background)
	count, err := q.Count(this.c)
	Check(this.c, err)
	
	var result bool
	if count > 1 {
		result = true
	} else {
		result = false
	}
	return result
}