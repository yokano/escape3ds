/**
 * ゲームの状態管理
 */
var State = Backbone.Model.extend({
	initialize: function() {
		this.set('currentScene', game.firstScene);
		this.set('itemList', new ItemList());
	}
});

/**
 * ゲーム
 */
var Game = Backbone.Model.extend({
	parse: function(data, options) {
		if(data.firstScene == '') {
			alert('最初のシーンが設定されていません');
		}
		this.set('name', data.name);
		this.set('description', data.description);
		this.set('itemList', new ItemList(data.itemList, {parse: true}));
		this.set('sceneList', new SceneList(data.sceneList, {parse: true}));
	}
});

/**
 * アイテム
 */
var Item = Backbone.Model.extend({
});

/**
 * アイテムリスト
 */
var ItemList = Backbone.Collection.extend({
	model: Item,
	parse: function(data, options) {
		_.each(data, function(val, key) {
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
});

/**
 * シーンリスト
 */
var SceneList = Backbone.Collection.extend({
	model: Scene,
	parse: function(data, options) {
		_.each(data, function(val, key) {
			this.add(new Scene({
				id: key,
				name: val.name,
				background: val.background,
				enter: val.enter,
				leave: val.leave,
				eventList: new EventList(val.eventList, {parse: true})
			}));
		}, this);
	}
});

/**
 * イベント
 */
var Event = Backbone.Model.extend({
	defaults: {
		visible: true  // 画面上に表示するかどうか
	},
	hide: function() {
		this.set('visible', false);
	},
	show: function() {
		this.set('visible', true);
	}
});

/**
 * イベントリスト
 */
var EventList = Backbone.Collection.extend({
	model: Event,
	parse: function(attr, options) {
		_.each(attr, function(val, key) {
			this.add(new Event({
				id: key,
				image: val.image,
				position: val.position,
				size: val.size,
				code: val.code
			}));
		}, this);
	}
});



