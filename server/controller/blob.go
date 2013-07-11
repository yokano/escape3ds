package controller

import (
	"appengine"
	"net/http"
	"fmt"
	. "server/lib"
	. "server/model"
)

// クライアントからアップロードされたファイルを blobstore に保存して blobkey を返す。
// Ajax で使う。method は POST。
func (this *Controller) Upload(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	
	file, fileHeader, err := r.FormFile("file")
	Check(c, err)
	
	model := NewModel(c)
	blobKey := model.AddBlob(file, fileHeader)
	
	fmt.Fprintf(w, `{"blobkey":"%s"}`, blobKey)
}

// blobstore から blogkey に関連付いたファイルをクライアントへ渡す。
// Ajax で使う。method は GET。
func (this *Controller) Download(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	
	blobKey := r.FormValue("blobkey")
	if blobKey == "" {
		c.Warningf("blobkey なしで download が実行されました")
		return
	}
	
	model := NewModel(c)
	contentType, bytes := model.GetBlob(blobKey)
	
	header := w.Header()
	header.Add("Content-Type", contentType)
	_, err := w.Write(bytes)
	Check(c, err)
}

func (this *Controller) Delete(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	blobKey := r.FormValue("blobkey")
	model := NewModel(c)
	model.DeleteBlob(blobKey)
	fmt.Fprintf(w, `{}`)
}