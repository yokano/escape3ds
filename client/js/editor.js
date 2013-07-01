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
		firstScene: null
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
	model: Scene,
	selected: null,
	select: function(cid) {
		this.selected = cid;
	},
	removed: function() {
		this.selected = null;
	},
	initialize: function() {
		this.on('select', this.select);
		this.on('remove', this.removed);
	}
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
		this.listenTo(this.collection, 'add remove', this.render);
		this.listenTo(this.collection, 'select', this.select);
		this.listenTo(game, 'change', this.gameHasChanged);
	},
	render: function() {
		this.$el.empty();
		this.collection.each(function(scene) {
			var sceneView = new SceneListItemView({
				model: scene,
				parent: this
			});
			this.$el.append(sceneView.render().el);
		}, this);
		return this;
	},
	select: function(cid) {
		var scene = this.collection.get(cid);
		var index = this.collection.indexOf(scene)
		this.$el.children().removeClass('select');
		this.$el.children().eq(index).addClass('select');
	},
	gameHasChanged: function() {
		this.$el.find('.is_first_scene').hide();
		var index = this.collection.indexOf(game.get('firstScene'));
		if(index == -1) {
			return;
		}
		this.$el.children().eq(index).find('.is_first_scene').show();
	}
});

/**
 * シーン一覧の１つ１つのli要素
 * @class
 * @extends Backbone.View
 * @member {SceneListView} parent 親要素への参照
 */
var SceneListItemView = Backbone.View.extend({
	tagName: 'li',
	template: _.template($('#scene_li_view_template').html()),
	parent: null,
	initialize: function(config) {
		this.parent = config.parent;
		this.listenTo(this.model, 'change', this.sceneHasChanged);
	},
	render: function() {
		this.$el.html(this.template(this.model.toJSON()));
		return this;
	},
	events: {
		'click': 'selectScene'
	},
	selectScene: function(event) {
		if(this.$el.hasClass('select')) {
			return;
		}
		this.parent.collection.trigger('select', this.model.cid)
	},
	sceneHasChanged: function() {
		this.$el.find('.scene_name').html(this.model.get('name'));
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
	model: null,
	initialize: function() {
		this.listenTo(sceneList, 'select', this.sceneHasChanged)
		this.listenTo(sceneList, 'remove', this.sceneHasRemoved)
	},
	render: function() {
		if(this.model == null) {
			this.$el.hide();
			return this;
		}
		this.$el.show();
		this.$el.html(this.template(this.model.toJSON()));
		if(game.get('firstScene') == this.model) {
			this.$el.find('#scene_info .is_first_scene').attr('checked', true);
		}
		return this;
	},
	events: {
		'change #change_scene_img': 'upload',
		'change #scene_info .scene_name': 'sceneNameHasChanged',
		'click #delete_scene': 'deleteButtonHasClicked',
		'click #copy_scene': 'copyButtonHasClicked',
		'click #scene_info .is_first_scene': 'firstSceneCheckboxHasClicked'
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
			}
		});
	},
	sceneHasChanged: function(cid) {
		this.model = sceneList.get(cid);
		this.render();
	},
	sceneNameHasChanged: function() {
		var name = this.$el.find('#scene_info .scene_name').val();
		this.model.set('name', name);
	},
	sceneHasRemoved: function() {
		this.model = null;
		this.render();
	},
	deleteButtonHasClicked: function() {
		var scene = sceneList.get(sceneList.selected);
		if(scene == null) {
			return;
		}
		if(!window.confirm(scene.get('name') + 'を削除しますか？')) {
			return;
		}
		sceneList.remove(scene);
	},
	copyButtonHasClicked: function() {
		var name = window.prompt('コピー先のシーン名を入力してください');
		if(name == '') {
			return;
		}
		var clone = this.model.clone();
		clone.set('name', name);
		sceneList.add(clone);
	},
	firstSceneCheckboxHasClicked: function() {
		var checked = this.$el.find('#scene_info .is_first_scene:checked').length;
		if(checked == 1) {
			game.set('firstScene', this.model);
		} else {
			game.set('firstScene', null);
		}
	}
});

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
