/**
 * エントリポイント
 * こちらは開発用のコード
 * デプロイ時にはミニファイしたものを使う
 */

// ファイル読み込み関数
var getFileURL = function(file, caller, callback) {
	var fileReader = new FileReader();
	fileReader.onload = function(data) {
		callback.call(caller, data.target.result);
	};
	fileReader.readAsDataURL(file);
};

//var game;
var sceneList;
var eventList;
var rootView;

rootView = new RootView();
$('body').html(rootView.render().el);
