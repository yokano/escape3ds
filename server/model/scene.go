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
	Name string `json:"name"`        // シーン名
	Background string `json:"background"`  // 背景画像のBlobkey、設定していなければ空文字
	Enter []byte `json:"-"`       // シーン開始時のイベント（データベース保存用バイナリ）
	RawEnter string `json:"enter" datastore:"-"`  // シーン開始時のイベント（クライアント用文字列）
	Leave string `json:"leave"`       // シーン終了時のイベント
	Sort int `json:"sort"`           // 並び順
	EventList map[string]*Event `json:"eventList" datastore:"-"` // シーン内のイベントリスト
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

// シーンをデータストアから取得する
func (this *Model) GetScene(encodedSceneKey string) *Scene {
	sceneKey := DecodeKey(this.c, encodedSceneKey)
	
	scene := new(Scene)
	err := datastore.Get(this.c, sceneKey, scene)
	Check(this.c, err)
	
	scene.EventList = this.GetEventList(encodedSceneKey)
	scene.RawEnter = string(scene.Enter)
	
	return scene
}

// 指定されたイベントが属しているシーンを取得する
func (this *Model) GetSceneFromEvent(encodedEventKey string) (string, *Scene) {
	eventKey, err := datastore.DecodeKey(encodedEventKey)
	Check(this.c, err)
	sceneKey := eventKey.Parent()
	encodedSceneKey := sceneKey.Encode()
	return encodedSceneKey, this.GetScene(encodedSceneKey)
}

// 引数として指定された id のシーンをデータストアから削除する。
func (this *Model) DeleteScene(encodedSceneKey string) {
	scene := this.GetScene(encodedSceneKey)

	// 内包するイベントを削除
	for key, _ := range scene.EventList {
		this.DeleteEvent(key)
	}

	// 他のシーンから画像が参照されていなければ画像を削除する
	if scene.Background != "" && !this.CheckDuplicateBackground(scene.Background) {
		this.DeleteBlob(scene.Background)
	}

	sceneKey := DecodeKey(this.c, encodedSceneKey)
	err := datastore.Delete(this.c, sceneKey)
	Check(this.c, err)
}

// id で指定されたシーンを新しいシーン scene で更新する。
func (this *Model) UpdateScene(encodedSceneKey string, scene *Scene) {
	sceneKey := DecodeKey(this.c, encodedSceneKey)
	sceneKey, err := datastore.Put(this.c, sceneKey, scene)
	Check(this.c, err)
}

// シーン内のイベントリストを取得する
func (this *Model) GetEvents(encodedSceneKey string) map[string]*Event {
	sceneKey, err := datastore.DecodeKey(encodedSceneKey)
	Check(this.c, err)
	
	query := datastore.NewQuery("Event").Ancestor(sceneKey)
	count, err := query.Count(this.c)
	Check(this.c, err)
	
	events := make([]*Event, count)
	eventKeys, err := query.GetAll(this.c, events)
	Check(this.c, err)
	
	encodedEventKeys := make([]string, count)
	for i := 0; i < count; i++ {
		encodedEventKeys[i] = eventKeys[i].Encode()
	}
	
	result := make(map[string]*Event, count)
	for i := 0; i < count; i++ {
		result[encodedEventKeys[i]] = events[i]
	}
	return result
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
		
		fmt.Fprintf(w, `{"sceneId":"%s"}`, encodedSceneKey)

	case "PUT":
		body := GetRequestBodyJSON(r)
		newScene := new(Scene)
		json.Unmarshal(body, newScene)
		
		encodedSceneKey := path[4]
		oldScene := this.GetScene(encodedSceneKey)
		
		if oldScene.Background != "" && oldScene.Background != newScene.Background && !this.CheckDuplicateBackground(oldScene.Background) {
			this.DeleteBlob(oldScene.Background)
		}
		
		this.UpdateScene(encodedSceneKey, newScene)
		
		fmt.Fprintf(w, `{}`)

	case "DELETE":
		encodedSceneKey := path[4]
		this.DeleteScene(encodedSceneKey)
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