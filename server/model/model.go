// datastore や blobstore などのサーバ上のデータを操作をまとめたもの。
// シーン、イベント、アイテムなどゲームに関わるデータや、
// 仮登録ユーザ、本登録ユーザ、セッション管理などのユーザに関するデータを扱う。
package model

import (
	"appengine"
)

// アプリケーションのデータ操作を行うオブジェクト
type Model struct {
	c appengine.Context
}

// モデルのインスタンスを作成する。
// モデルを使う前にかならず実行すること。
func NewModel(c appengine.Context) *Model {
	model := new(Model)
	model.c = c
	return model
}
