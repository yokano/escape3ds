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

var game;
var sceneList;
var eventList;
var rootView;

game = new Game({
	name: '誕生日の脱出劇',
	description: '誕生日会の翌日に目が覚めると見知らぬ部屋に閉じ込められていた'
});

game.get('sceneList').add([
	{name: '押入れ１'},
	{name: '押入れ２'}
]);

rootView = new RootView();
$('body').html(rootView.render().el);
