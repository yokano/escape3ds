var game = new Game({}, data);

var state = new State({
	currentScene: game.firstScene,
	debug: true
});

var rootView = new RootView({
	model: state
});
$('body').append(rootView.render().el);
