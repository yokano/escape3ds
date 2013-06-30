$(function() {

/**
 * モデルの定義
 */
var Game = Backbone.Model.extend({
	defaults: {
		name: '',
		description: '',
		thumbnail: '',
		userKey: '',
		firstScene: ''
	}
});

var Scene = Backbone.Model.extend({
	defaults: {
		name: '',
		gameKey: '',
		background: 'black.png',
		enter: '',
		leave: ''
	}
});

var Event = Backbone.Model.extend({
	defaults: {
		name: '',
		image: '',
		code: '',
		position: '',
		size: '',
		sceneKey: ''
	}
});

var Item = Backbone.Model.extend({
	defaults: {
		name: '',
		image: '',
		gameKey: ''
	}
});

var SceneList = Backbone.Collection.extend({
	model: Scene
});

var EventList = Backbone.Collection.extend({
	model: Event
});

var ItemList = Backbone.Collection.extend({
	model: Item
});


/**
 * ビューの定義
 *
 * rootView
 *   HeaderView
 *   SceneEditorView
 *     SceneListView
 *       SceneListItemView
 *     SceneView
 *     EventView
 */
var RootView = Backbone.View.extend({
	tagName: 'div',
	render: function() {
		var headerView = new HeaderView({model: game});
		var sceneEditorView = new SceneEditorView()
		this.$el.append(
			headerView.render().el,
			sceneEditorView.render().el
		);
		return this;
	}
});

/**
 * ヘッダービュー
 * 画面上に表示されるヘッダ部分のビュー
 * 表示するモデル : game
 * @class
 * @extends Backbone.View
 */
var HeaderView = Backbone.View.extend({
	tagName: 'header',
	template: _.template($('#header_view_template').html()),
	render: function() {
		this.$el.html(this.template(this.model.toJSON()));
		return this;
	},
	events: {
		'click #scene_mode' : 'changeToSceneMode',
		'click #item_mode': 'changeToItemMode',
		'click #back' : 'backToGameList',
		'click #save' : 'save'
	},
	changeToSceneMode: function() {
		console.log('シーン管理画面を表示');
	},
	changeToItemMode: function() {
		console.log('アイテム管理画面を表示');
	},
	backToGameList: function() {
		console.log('ゲームリストへ戻る');
	},
	save: function() {
		console.log('ゲームデータを保存する');
	}
});

/**
 * シーンエディタビュー
 * シーンの編集を行うビュー
 * ヘッダの下に表示される
 * 左にシーンリスト、中央にシーン、右にイベントエディタを表示
 * @class
 * @extends Backbone.View
 */
var SceneEditorView = Backbone.View.extend({
	tagname: 'section',
	render: function() {
		var sceneListView = new SceneListView({collection: sceneList});
		var sceneView = new SceneView();
		var eventView = new EventView();
		
		$('<button id="add_scene">シーンを追加</button>').appendTo(this.$el);
		this.$el.append(
			sceneListView.render().el,
			sceneView.render().el,
			eventView.render().el
		);
		return this;
	},
	events: {
		'click #add_scene': 'addScene'
	},
	addScene: function() {
		var name = window.prompt('シーン名');
		if(name == '') {
			return;
		}
		sceneList.add({
			name: name
		});
	}
});

/**
 * シーンリストビュー
 * @class
 * @extends Backbone.View
 */
var SceneListView = Backbone.View.extend({
	tagName: 'ul',
	id: 'scene_list',
	initialize: function() {
		this.listenTo(this.collection, 'add', this.render);
	},
	render: function() {
		this.$el.empty();
		this.collection.each(function(scene) {
			var sceneView = new SceneListItemView({model: scene});
			this.$el.append(sceneView.render().el);
		}, this);
		return this;
	}
});

/**
 * シーン一覧の１つ１つのli要素
 * @class
 * @extends Backbone.View
 */
var SceneListItemView = Backbone.View.extend({
	tagName: 'li',
	template: _.template($('#scene_li_view_template').html()),
	render: function() {
		this.$el.html(this.template(this.model.toJSON()));
		return this;
	}
});

/**
 * シーンの設定ビュー
 * 画面中央
 * @class
 * @extends Backbone.View
 */
var SceneView = Backbone.View.extend({
	tagName: 'div',
	template: _.template($('#scene_view_template').html()),
	render: function() {
		this.$el.html(this.template());
		return this;
	},
	events: {
		'change #change_scene_img': 'upload'
	},
	upload: function(data) {
		var form = $('#change_scene_img_form').get()[0];
		var formData = new FormData(form);
		$.ajax('/upload', {
			method: 'POST',
			contentType: false,
			processData: false,
			data: formData,
			dataType: 'json',
			error: function() {
				console.log('error');
			},
			success: function(data) {
				console.log('blobkey', data.blobkey);
				download(data.blobkey);
			}
		});
	}
});

var download = function(blobKey) {
	$.ajax('/donwload', {
		method: 'GET',
		data: {
			blobKey: blobKey
		},
		dataType: 'json',
		error: function() {
			console.log('error');
		},
		success: function(data) {
			console.log(data);
		}
	});
};

/**
 * イベント編集ビュー
 * 画面右側
 * @class
 * @extends Backbone.View
 */
var EventView = Backbone.View.extend({
	tagName: 'div',
	template: _.template($('#event_view_template').html()),
	render: function() {
		this.$el.html(this.template());
		return this;
	}
});


/**
 * エントリポイント
 */
var game = new Game({
	name: '誕生日の脱出劇',
	description: '誕生日会の翌日に目が覚めると見知らぬ部屋に閉じ込められていた',
	thumbnail: '',
	userKey: '',
	firstScene: ''
});

var sceneList = new SceneList([
	{name: '押入れ１', background: 'armoire1.png'},
	{name: '押入れ２', background: 'armoire2.png'}
]);

var rootView = new RootView();
$('body').html(rootView.render().el);

});
