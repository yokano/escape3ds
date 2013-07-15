/**
 * エントリポイント
 * こちらは開発用のコード
 * デプロイ時にはミニファイしたものを使う
 * @author y.okano
 * @file
 */

// ビューを画面に表示
var rootView = new RootView();
$('body').html(rootView.render().el);
