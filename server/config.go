/**
 * 設定ファイル
 */
package escape3ds
var config = map[string]string {
	"consumer_key": "j44BdfM2pzUISfPc5s3YZw",
	"url": "https://api.twitter.com/oauth/request_token",
	"consumer_secret": "vXB86muAkDhyjFSCQACpqBjVQnmVLKxVroDdGMwZEH0",
	"facebook_client_id": "305513316248719",
	"facebook_client_secret": "3f6789addbba0bd9413af15b3695c8df",
	"interimMailBody": `

%s 様

このメールは ESCAPE 3DS へご登録いただいたメールアドレスへ自動送信しております。
こちら Web サイトに心当たりの無い方は、メールを無視してください。

ESCPAE 3DS
http://escape-3ds.appspot.com/


新規登録なさったご本人の場合は下の URL をクリックして本登録を完了してください。


こちらの URL をクリックして本登録
http://escape-3ds.appspot.com/registration?key=%s


お問い合わせはこちらのメールアドレスではなく yuta.okano@gmail.com へお願い致します。
この度は ESCAPE 3DS へご登録いただきましてありがとうございました。

---------------------------------
- ESCPAE 3DS
- 開発者：岡野雄太(Yuta Okano)
- 連絡先：yuta.okano@gmail.com
---------------------------------

`,
}