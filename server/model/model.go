/**
 * データモデルの定義
 * @file
 */
package escape3ds

import (
	"appengine"
)

/**
 * モデル
 * @class
 * @property {appengine.Context} c コンテキスト
 */
type Model struct {
	c appengine.Context
}

/**
 * モデルの作成
 * @function
 * @param {appengine.Context} c コンテキスト
 * @returns {*Model} モデル
 */
func NewModel(c appengine.Context) *Model {
	model := new(Model)
	model.c = c
	return model
}
