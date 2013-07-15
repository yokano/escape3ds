package model

import (
	. "server/lib"

	"appengine/datastore"
	"net/http"
	"encoding/json"
	"fmt"
)

// アイテム型
type Item struct {
	Name string
	Img string
	HasFirst bool
}

// アイテムの作成
func NewItem(name string, img string, hasFirst bool) *Item {
	item := new(Item)
	item.Name = name
	item.Img = img
	item.HasFirst = hasFirst
	return item
}

// アイテムの追加。
// 追加したいアイテムと、追加するゲームキーを渡す。
// 追加したアイテムのエンコード済みキーを返す。
func (this *Model) AddItem(item *Item, encodedGameKey string) string {
	gameKey, err := datastore.DecodeKey(encodedGameKey)
	Check(this.c, err)
	
	itemKey := datastore.NewIncompleteKey(this.c, "Item", gameKey)
	itemKey, err = datastore.Put(this.c, itemKey, item)
	Check(this.c, err)
	
	encodedItemKey := itemKey.Encode()
	return encodedItemKey
}

// アイテム一覧の取得。
func (this *Model) GetItems(encodedGameKey string) map[string]*Item {
	gameKey, err := datastore.DecodeKey(encodedGameKey)
	Check(this.c, err)
	
	query := datastore.NewQuery("Item").Ancestor(gameKey)
	count, err := query.Count(this.c)
	Check(this.c, err)

	iterator := query.Run(this.c)
	items := make(map[string]*Item, count)
	
	for ;; {
		item := new(Item)
		itemKey, err := iterator.Next(item)
		if err != nil {
			break
		}
		encodedItemKey := itemKey.Encode()
		items[encodedItemKey] = item
	}
	
	return items
}

// アイテムの同期
func (this *Model) SyncItem(w http.ResponseWriter, r *http.Request, path []string) {
	switch r.Method {
	case "POST":
		body := GetRequestBodyJSON(r)
		item := new(Item)
		json.Unmarshal(body, item)
		
		gameKey, err := datastore.DecodeKey(path[3])
		Check(this.c, err)
		
		itemKey := datastore.NewIncompleteKey(this.c, "Item", gameKey)
		itemKey, err = datastore.Put(this.c, itemKey, item)
		Check(this.c, err)
		
		encodedItemKey := itemKey.Encode()
		
		fmt.Fprintf(w, `{"itemKey":"%s"}`, encodedItemKey)
		
	case "PUT":
		this.c.Debugf("PUT")
	case "GET":
		this.c.Debugf("GET")
	case "DELETE":
		this.c.Debugf("DELETE")
	}
}