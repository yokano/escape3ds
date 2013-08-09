/**
 * ゲームの状態管理
 */
var State = Backbone.Model.extend({
	defaults: {
		currentScene: null,  // 最初に表示するシーンオブジェクト
		itemList: null,  // 所持しているアイテムリスト
		busy: false  // メッセージのページ送り待ちなら true
	},
	initialize: function() {
		this.set('currentScene', game.get('firstScene'));

		var firstItems = {};
		_.each(data.itemList, function(val, key) {
			if(val.hasFirst) {
				firstItems[key] = val;
			}
		});
		this.set('itemList', new ItemList(null, firstItems));
	},
	changeScene: function(id) {
		this.set('currentScene', game.get('sceneList').get(id));
	},
	addItem: function(id) {
		var item = game.get('itemList').get(id)
		this.get('itemList').add(item);
	},
	removeItem: function(id) {
		this.get('itemList').remove(id);
	},
	getCurrentItem: function() {
		var itemList = this.get('itemList');
		var currentItem = itemList.at(itemList.cursor);
		if(currentItem == undefined) {
			currentItem = null;
		}
		return currentItem;
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
		firstScene: null,  // 最初のシーン
		message: null  // メッセージ
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
		this.set('message', new Message());
	}
});

/**
 * メッセージ
 */
var Message = Backbone.Model.extend({
	defaults: {
		current: 0,  // 現在表示中のページ番号
		queue: []  // メッセージキュー
	},
	show: function(messageQueue) {
		this.set('current', 0);
		this.set('queue', messageQueue);
		if(this.get('queue').length > 1) {
			state.set('busy', true);
		}
	},
	next: function() {
		this.set('current', this.get('current') + 1);
		if(this.get('current') >= this.get('queue').length - 1) {
			state.set('busy', false);
		}
	}
});

/**
 * アイテム
 */
var Item = Backbone.Model.extend({
	defaults: {
		name: '',
		hasFirst: '',
		img: '',
		selected: false
	}
});

/**
 * アイテムリスト
 */
var ItemList = Backbone.Collection.extend({
	model: Item,
	cursor: -1,  // 現在選択中のアイテムのindex, -1 は未選択状態
	initialize: function(attr, options) {
		_.each(options, function(val, key) {
			this.add(new Item({
				id: key,
				name: val.name,
				hasFirst: val.has_first,
				img: val.img
			}));
		}, this);
	},
	next: function() {
		if(this.length == 0) {
			return;
		}
		
		if(this.cursor == -1) {
			this.cursor = 0;
		} else if(this.cursor >= this.length - 1) {
			this.at(this.cursor).set('selected', false);
			this.cursor = 0;
		} else {
			this.at(this.cursor).set('selected', false);
			this.cursor++;
		}
		
		this.at(this.cursor).set('selected', true);
	},
	prev: function() {
		if(this.length == 0) {
			return;
		}
		
		if(this.cursor == -1) {
			this.cursor = this.length - 1;
		} else if(this.cursor <= 0) {
			this.at(this.cursor).set('selected', false);
			this.cursor = this.length - 1;
		} else {
			this.at(this.cursor).set('selected', false);
			this.cursor--;
		}
		
		this.at(this.cursor).set('selected', true);
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
		console.log('remove');
	},
	execute: function() {
		_.each(this.get('code'), function(method) {
			switch(method.type) {
			case 'changeScene': {
				state.changeScene(method.attr);
				break;
			}
			case 'addItem': {
				state.addItem(method.attr);
				break;
			}
			case 'removeItem': {
				state.removeItem(method.attr);
				break;
			}
			case 'remove': {
				this.remove();
				break;
			}
			case 'hide': {
				this.hide();
				break;
			}
			case 'show': {
				this.show();
				break;
			}
			default: {
				console.log('不明なイベント内容が実行されました', method);
			}
			}
		}, this);
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
				code: JSON.parse(val.rawcode)
			});
		}, this);
		
		this.on('change:removed', this.removeEvent);
	},
	removeEvent: function(event) {
		// イベントが削除された時の処理
	}
});
