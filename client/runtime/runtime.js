var game = new Game({}, data);

var state = new State({
	currentScene: game.firstScene,
	debug: true
});

var rootView = new RootView({
	model: state
});
$('body').append(rootView.render().el);

// スクロールを促す
game.get('message').show(['十字キーの↓を押して画面を合わせてください']);

window.onscroll = function() {
	alert('scroll');
};