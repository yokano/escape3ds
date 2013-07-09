// エントリポイント
// ESCAPE 3DS は Nintendo 3DS の Web ブラウザ上で動作する
// 脱出ゲームを開発するための Web アプリケーションです
package escape3ds

import (
	. "server/controller"
)

// コントローラの生成と処理の振り分け
func init() {
	controller := NewController()
	controller.Handle()
}