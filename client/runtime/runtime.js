var game = new Game({}, data);

var state = new State({
	currentScene: game.firstScene,
	debug: true
});

//state.get('itemList').add([
//	new Item(game.item_list.hammer, {parse: true}),
//	new Item(game.item_list.dish, {parse: true}),
//]);

var rootView = new RootView({
	model: state
});
$('body').append(rootView.render().el);
