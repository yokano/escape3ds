package model

import (
	"mime/multipart"
	"appengine"
	"appengine/datastore"
	"appengine/blobstore"
	. "server/lib"
)

/**
 * ファイルを blobstore に保存して blobkey を返す
 * @method
 * @memberof Model
 * @param {multipart.File} file 保存するファイル
 * @param {*multipart.FileHeader} fileHeader 保存するファイルのヘッダ
 */
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

/**
 * blobstore に保存したファイルを読み出す
 * @method
 * @memberof Model
 * @param {string} key 読みだすファイルの blobkey
 * @returns {string} ファイルの Content-Type
 * @returns {[]byte} ファイルバイナリデータ
 */
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