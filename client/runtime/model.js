/**
 * ゲームの状態管理
 */
var State = Backbone.Model.extend({
	defaults: {
		currentScene: null,  // 最初に表示するシーンオブジェクト
		itemList: null,  // 所持しているアイテムリスト
	},
	initialize: function() {
		this.set('currentScene', game.get('firstScene'));
		this.set('itemList', new ItemList());
	},
	
	/**
	 * 別のシーンへ移動する
	 * @param {String} id 移動先シーンのid
	 */
	changeScene: function(id) {
		this.set('currentScene', game.get('sceneList').get(id));
	}
});

/**
 * ゲーム
 */
var Game = Backbone.Model.extend({
	defaults: {
		name: '',  // ゲーム名
		description: '',  // ゲームの説明
		itemList: null,  // ゲームに存在するすべてのアイテム
		sceneList: null,  // ゲームに存在するすべてのシーン
		firstScene: null  // 最初のシーン
	},
	initialize: function(attr, options) {
		if(options.firstScene == '') {
			alert('最初のシーンが設定されていません');
		}
		this.set('name', options.name);
		this.set('description', options.description);
		this.set('itemList', new ItemList(null, options.itemList));
		this.set('sceneList', new SceneList(null, options.sceneList));
		this.set('firstScene', this.get('sceneList').get(options.firstScene));
	}
});

/**
 * アイテム
 */
var Item = Backbone.Model.extend({
	defaults: {
		name: '',
		hasFirst: '',
		img: ''
	}
});

/**
 * アイテムリスト
 */
var ItemList = Backbone.Collection.extend({
	model: Item,
	initialize: function(attr, options) {
		_.each(options, function(val, key) {
			this.add(new Item({
				id: key,
				name: val.name,
				hasFirst: val.has_first,
				img: val.img
			}));
		}, this);
	}
});

/**
 * シーン
 */
var Scene = Backbone.Model.extend({
	defaults: {
		name: '',
		background: '',
		enter: '',
		leave: '',
		eventList: null
	}
});

/**
 * シーンリスト
 */
var SceneList = Backbone.Collection.extend({
	model: Scene,
	initialize: function(attr, options) {
		_.each(options, function(val, key) {
			this.add({
				id: key,
				name: val.name,
				background: val.background,
				enter: val.enter,
				leave: val.leave,
				eventList: new EventList(null, val.eventList)
			});
		}, this);
	}
});

/**
 * イベント
 */
var Event = Backbone.Model.extend({
	defaults: {
		visible: true,  // 画面上に表示するかどうか
		removed: false,
	},
	hide: function() {
		this.set('visible', false);
	},
	show: function() {
		this.set('visible', true);
	},
	remove: function() {
		this.set('removed', true);
	}
});

/**
 * イベントリスト
 */
var EventList = Backbone.Collection.extend({
	model: Event,
	initialize: function(attr, options) {
		_.each(options, function(val, key) {
			this.add({
				id: key,
				image: val.image,
				position: val.position,
				size: val.size,
				code: val.code
			});
		}, this);
		
		this.on('change:removed', this.removeEvent);
	},
	removeEvent: function(event) {
		// イベントが削除された時の処理
	}
});



