// アプリケーションのデータ操作を行う
package model

import (
	"appengine"
)

// アプリケーションのデータ操作を行うモデル
// c: アプリケーションコンテキスト
type Model struct {
	c appengine.Context
}

// モデルの新規作成作成
// c アプリケーションコンテキスト
func NewModel(c appengine.Context) *Model {
	model := new(Model)
	model.c = c
	return model
}
