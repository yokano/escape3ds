var blockList = new BlockList(code);

var rootView = new RootView();
$('body').html(rootView.render().el);
