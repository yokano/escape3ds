package model
import(
	"appengine/datastore"
	"net/http"
	"encoding/json"
	"fmt"
	
	. "server/lib"
)

// シーンに含まれるイベントオブジェクト
type Event struct {
	Name string `json:"name"`     // イベント名
	Image string `json:"image"`    // 画像
	Code string `json:"code"`     // クリックされた時のコード
	Position []int `json:"position"`  // 位置
	Size []int `json:"size"`      // 大きさ
	Color string `json:"color"`    // エディタ上で領域に表示される色
}

// イベントオブジェクトの作成
func (this *Model) NewEvent(name string, image string, code string, position []int, size []int, color string) *Event {
	event := new(Event)
	event.Name = name
	event.Image = image
	event.Code = code
	event.Position = position
	event.Size = size
	event.Color = color
	return event
}

// イベントの追加。エンコード済みのキーを返す。
// 引数としてイベントが蔵するシーンキーを渡す。
func (this *Model) AddEvent(event *Event, encodedSceneKey string) string {
	sceneKey, err := datastore.DecodeKey(encodedSceneKey)
	Check(this.c, err)
	
	eventKey := datastore.NewIncompleteKey(this.c, "Event", sceneKey)
	eventKey, err = datastore.Put(this.c, eventKey, event)
	Check(this.c, err)
	
	encodedEventKey := eventKey.Encode()
	return encodedEventKey
}

// イベントを取得
func (this *Model) GetEvent(encodedEventKey string) *Event {
	eventKey := DecodeKey(this.c, encodedEventKey)
	event := new(Event)
	err := datastore.Get(this.c, eventKey, event)
	Check(this.c, err)
	return event
}

// イベントを更新する
func (this *Model) UpdateEvent(event *Event, encodedEventKey string) {
	eventKey, err := datastore.DecodeKey(encodedEventKey)
	Check(this.c, err)
	
	_, err = datastore.Put(this.c, eventKey, event)
	Check(this.c, err)
}

// イベントの削除
func (this *Model) DeleteEvent(encodedEventKey string) {
	eventKey, err := datastore.DecodeKey(encodedEventKey)
	Check(this.c, err)
	
	err = datastore.Delete(this.c, eventKey)
	Check(this.c, err)
}

// 指定されたシーンに属するイベント一覧を返す
func (this *Model) GetEventList(encodedSceneKey string) map[string]*Event {
	sceneKey := DecodeKey(this.c, encodedSceneKey)
	query := datastore.NewQuery("Event").Ancestor(sceneKey)
	count, err := query.Count(this.c)
	Check(this.c, err)
	
	eventList := make(map[string]*Event, count)
	iterator := query.Run(this.c)
	for i := 0; i < count; i++ {
		eventKey, err := iterator.Next(nil)
		Check(this.c, err)
		
		encodedEventKey := eventKey.Encode()
		eventList[encodedEventKey] = this.GetEvent(encodedEventKey)
	}
	
	return eventList
}

// イベントコードの更新
func (this *Model) UpdateEventCode(id string, code string) {
	event := this.GetEvent(id)
	event.Code = code
	this.UpdateEvent(event, id)
	this.c.Debugf("UPDATE CODE")
	this.c.Debugf("ID: %s", id)
	this.c.Debugf("CODE: %s", code)
}

// イベントの同期
func (this *Model) SyncEvent(w http.ResponseWriter, r *http.Request, path []string) {
	switch r.Method {
	case "POST":
		this.c.Debugf("CREATE EVENT")
		body := GetRequestBodyJSON(r)
		event := new(Event)
		json.Unmarshal(body, event)
		
		sceneKey, err := datastore.DecodeKey(path[3])
		Check(this.c, err)
		
		eventKey := datastore.NewIncompleteKey(this.c, "Event", sceneKey)
		eventKey, err = datastore.Put(this.c, eventKey, event)
		Check(this.c, err)
		
		encodedEventKey := eventKey.Encode()
		
		fmt.Fprintf(w, `{"id":"%s"}`, encodedEventKey)
		
	case "PUT":
		this.c.Debugf("UPDATE EVENT")
		body := GetRequestBodyJSON(r)
		event := new(Event)
		json.Unmarshal(body, event)
		
		eventKey, err := datastore.DecodeKey(path[4])
		Check(this.c, err)

		eventKey, err = datastore.Put(this.c, eventKey, event)
		Check(this.c, err)
		
		fmt.Fprintf(w, `{}`)
		
	case "GET":
		this.c.Debugf("GET")
		
	case "DELETE":
		encodedEventKey := path[4]
		eventKey, err := datastore.DecodeKey(encodedEventKey)
		Check(this.c, err)
		
		event := new(Event)
		err = datastore.Get(this.c, eventKey, event)
		Check(this.c, err)
		
//		this.DeleteBlob(event.Img)
		
		err = datastore.Delete(this.c, eventKey)
		Check(this.c, err)

		fmt.Fprintf(w, `{}`)
	}
}