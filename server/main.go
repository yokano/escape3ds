/**
 * エントリポイント
 * URLパターンから該当する処理へ振り分ける
 * 処理は controller.go に記載されている
 * @file
 */
package escape3ds

import (
	. "server/controller"
)


/**
 * URLから処理を振り分ける
 * @function
 */
func init() {
	controller := NewController()
	controller.Handle()
}