// ESCAPE 3DS は Nintendo 3DS の Web ブラウザ上で動作する
// 脱出ゲームを開発するための Web アプリケーションです
//
// ディレクトリ構成
//   model: データ操作
//   view: 画面表示
//   controller: URLパターンによる処理の振り分け
//   lib: 汎用ライブラリ
//   config: アプリの設定
//
package main

import (
	. "server/controller"
)

// エントリポイント
func init() {
	controller := NewController()
	controller.Handle()
}