var blockList;

// データが破損していたら新規作成
try {
	blockList = new BlockList(code);
} catch(e) {
	console.log(e);
	blockList = new BlockList();
}

var rootView = new RootView();
$('body').html(rootView.render().el);
