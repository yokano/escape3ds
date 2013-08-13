package model

import (
	"mime/multipart"
	"appengine"
	"appengine/datastore"
	"appengine/blobstore"
	. "server/lib"
)

// ファイルを blobstore に保存して blobkey を返す。
// エディタから画像ファイルをアップロードするときに使う。
// 引数としてファイルと、ファイルの種類などを含むファイルヘッダを渡す。
func (this *Model) AddBlob(file multipart.File, fileHeader *multipart.FileHeader) string {
	// ファイルデータを読み出す
	fileData := make([]byte, 256 * 1024)
	byteNum, err := file.Read(fileData)
	Check(this.c, err)
	fileData = fileData[:byteNum]
	
	// blobstore に blob を作成
	mimeType := fileHeader.Header.Get("Content-Type")
	writer, err := blobstore.Create(this.c, mimeType)
	Check(this.c, err)
	
	// blob にファイルデータを書き込む
	writedByteNum, err := writer.Write(fileData)
	Check(this.c, err)
	if writedByteNum != byteNum {
		this.c.Errorf("blob の書き込みに失敗しました")
		return ""
	}
	err = writer.Close()
	Check(this.c, err)
	
	// blob key を取得
	key, err := writer.Key()
	Check(this.c, err)
	
	return string(key)
}

// blobstore に保存したファイルを読み出す。
// 引数として読みだすファイルの blobkey を渡す。
// ファイルタイプとバイナリデータが返される。
func (this *Model) GetBlob(key string) (string, []byte) {
	blobKey := appengine.BlobKey(key)
	
	// __BlobFileIndex__ から blobKey で Key Name (string id) を取得
	query := datastore.NewQuery("__BlobFileIndex__").Filter("blob_key =", blobKey)
	iterator := query.Run(this.c)
	blobFileIndexKey, err := iterator.Next(nil)
	Check(this.c, err)
	
	// __Blobinfo__ から Key Name で size を取得
	query = datastore.NewQuery("__BlobInfo__").Filter("creation_handle =", blobFileIndexKey.StringID())
	iterator = query.Run(this.c)
	blobInfo := new(blobstore.BlobInfo)
	_, err = iterator.Next(blobInfo)
	Check(this.c, err)
	size := blobInfo.Size
	
	// size バイトだけ reader からバイトを読み出す
	reader := blobstore.NewReader(this.c, blobKey)
	bytes := make([]byte, size)
	readed, err := reader.Read(bytes)
	Check(this.c, err)
	if(int64(readed) != size) {
		this.c.Warningf("読み込んだファイルサイズとメタ情報のファイルサイズが異なります")
	}
	
	return blobInfo.ContentType, bytes
}

// blogstore に保存されているフィアルを削除する。
// 引数として削除する blobkey を渡す。
func (this *Model) DeleteBlob(key string) {
	blobKey := appengine.BlobKey(key)
	err := blobstore.Delete(this.c, blobKey)
	Check(this.c, err)
	if err != nil {
		this.c.Warningf("存在しないblobを削除しようとしました")
	}
}

// すべてのblobを削除する。管理者専用。
func (this *Model) ClearBlob() {
	// __BlobInfo__ のクリア
	q := datastore.NewQuery("__BlobInfo__").KeysOnly()
	
	blobInfoKeys, err := q.GetAll(this.c, nil)
	Check(this.c, err)

	blobKeys := make([]appengine.BlobKey, len(blobInfoKeys))
	for i := 0; i < len(blobInfoKeys); i++ {
		blobKeys[i] = appengine.BlobKey(blobInfoKeys[i].StringID())
	}
	
	err = blobstore.DeleteMulti(this.c, blobKeys)
	Check(this.c, err)
	
	// __BlobFileIndex__ のクリア
	q = datastore.NewQuery("__BlobFileIndex__").KeysOnly()
	indexKeys, err := q.GetAll(this.c, nil)
	err = datastore.DeleteMulti(this.c, indexKeys)
	Check(this.c, err)
}