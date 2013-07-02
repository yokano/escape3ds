/**
 * エントリポイント
 * こちらは開発用のコード
 * デプロイ時にはミニファイしたものを使う
 */
var game;
var sceneList;
var eventList;
var rootView;

game = new Game({
	name: '誕生日の脱出劇',
	description: '誕生日会の翌日に目が覚めると見知らぬ部屋に閉じ込められていた',
	thumbnail: '',
	userKey: '',
	firstScene: ''
});

sceneList = new SceneList([
	{name: '押入れ１'},
	{name: '押入れ２'}
]);

rootView = new RootView();
$('body').html(rootView.render().el);
