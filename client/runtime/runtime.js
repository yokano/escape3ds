var g = new Game(game, {parse: true});

var state = new State();
state.currentScene = game.sceneList[state.firstScene]

//state.get('itemList').add([
//	new Item(game.item_list.hammer, {parse: true}),
//	new Item(game.item_list.dish, {parse: true}),
//]);

var rootView = new RootView({
	model: g
});
$('body').append(rootView.render().el);