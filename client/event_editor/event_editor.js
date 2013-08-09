var blockList;

// データが破損していたら新規作成
try {
	blockList = new BlockList(null, code, true);
} catch(e) {
	console.log(e);
	blockList = new BlockList();
}

var rootView = new RootView();
$('body').html(rootView.render().el);


/*
var V1 = Backbone.View.extend({
	initialize: function() {
		console.log('V! initialize');
	}
});

var V2 = V1.extend({
	constructor: function() {
		V1.call(this);
		console.log('V2 initialize');
	}
});

var v2 = new V2();
*/