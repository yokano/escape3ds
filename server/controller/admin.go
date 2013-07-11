package controller

import (
	"net/http"
	"appengine"
	"fmt"
	"encoding/json"
	. "server/model"
	. "server/view"
	. "server/lib"
)

// デバッグツールの表示
func (this *Controller) Debug(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	view := NewView(c, w)
	view.Debug()
}

// 仮登録ユーザ一覧の取得。Ajax で呼び出す
func (this *Controller) GetInterimUsers(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	model := NewModel(c)
	interimUsers := model.GetInterimUsers()
	
	// キーと名前だけを返す
	result := make(map[string]string, len(interimUsers))
	for key, val := range interimUsers {
		result[key] = val.Name
	}
	
	bytes, err := json.Marshal(result)
	Check(c, err)
	fmt.Fprintf(w, "%s", bytes)
}

// ユーザ一覧の取得。Ajax で呼び出す
func (this *Controller) GetUsers(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	model := NewModel(c)
	users := model.GetAllUser()
	
	result := make(map[string]string, len(users))
	for key, val := range users {
		result[key] = val.Name
	}
	
	bytes, err := json.Marshal(result)
	Check(c, err)
	fmt.Fprintf(w, "%s", bytes)
}

// blobsotre のクリア。管理者専用。
func (this *Controller) ClearBlob(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	model := NewModel(c)
	model.ClearBlob()
	fmt.Fprintf(w, "{}")
}