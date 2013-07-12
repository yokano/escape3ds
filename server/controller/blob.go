package controller

import (
	"appengine"
	"appengine/blobstore"
	"net/http"
	"fmt"
	. "server/lib"
	. "server/model"
)

// クライアントからファイルがアップロードされたらblobkeyを返す
func (this *Controller) Uploaded(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	blobs, _, err := blobstore.ParseUpload(r)
	Check(c, err)
	fmt.Fprintf(w, `{"blobkey":"%s"}`, blobs["file"][0].BlobKey)
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
	
	blobstore.Send(w, appengine.BlobKey(blobKey))
}

// blobstore から指定された blobkey の Blob を削除する。
// 既に設定されている画像を変更した時に、以前の画像を削除するために使う。
func (this *Controller) Delete(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	blobKey := r.FormValue("blobkey")
	model := NewModel(c)
	model.DeleteBlob(blobKey)
	fmt.Fprintf(w, `{}`)
}

// blobstore にファイルをアップロードするためのURLを返す。
func (this *Controller) GetURL(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	url, err := blobstore.UploadURL(c, "/uploaded", nil)
	Check(c, err)
	
	fmt.Fprintf(w, `{"url":"%s"}`, url.String())
}