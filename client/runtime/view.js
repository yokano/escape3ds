/**
 * ルートビュー
 * model: State
 */
var RootView = Backbone.View.extend({
	id: 'root',
	tagName: 'div',
	render: function() {
		var upperView = new UpperView({
			model: this.model
		});
		this.$el.append(upperView.render().el);
		
		var sceneView = new SceneView({
			model: this.model
		});
		this.$el.append(sceneView.render().el);
				
		return this;
	}
});

/**
 * 上画面
 * model: State
 */
var UpperView = Backbone.View.extend({
	id: 'upper',
	tagName: 'div',
	render: function() {
		var itemListView = new ItemListView({
			collection: this.model.get('itemList')
		});
		this.$el.append(itemListView.render().el);
		
		var messageView = new MessageView({
			model: game.get('message')
		});
		this.$el.append(messageView.render().el);
		
		return this;
	}
});

/**
 * アイテムリスト
 * collection: ItemList
 */
var ItemListView = Backbone.View.extend({
	id: 'item_list',
	tagName: 'ul',
	render: function() {
		this.$el.empty();
		
		// 手持ちのアイテムを表示
		state.get('itemList').each(function(item) {
			var itemView = new ItemView({
				model: item
			});
			this.$el.append(itemView.render().el);
		}, this);
		
		// アイテムが10個以下なら空欄を表示
		for(var i = 0; i < 10 - state.get('itemList').length; i++) {
			this.$el.append($('<li class="item"></li>'));
		}
		
		return this;
	},
	initialize: function() {
		var view = this;
		$(document).on('keydown', function(event) {
			view.keyHasDown.call(view, event);
		});
		
		this.listenTo(this.collection, 'add remove', this.render);
	},
	keyHasDown: function(event) {
		var LEFT = 37;
		var RIGHT = 39;
		var keyCode = event.originalEvent.keyCode;
		if(keyCode == LEFT) {
			this.collection.prev();
		} else if(keyCode == RIGHT) {
			this.collection.next();
		}
	}
});

/**
 * 各アイテム
 */
var ItemView = Backbone.View.extend({
	className: 'item',
	tagName: 'li',
	render: function() {
		this.$el.append('<img src="/download?blobkey=' + this.model.get('img') + '">');
		return this;
	},
	initialize: function() {
		this.listenTo(this.model, 'change:selected', this.itemHasSelected);
	},
	itemHasSelected: function() {
		if(this.model.get('selected')) {
			this.$el.addClass('selected');
		} else {
			this.$el.removeClass('selected');
		}
	}
});

/**
 * メッセージ
 */
var MessageView = Backbone.View.extend({
	id: 'message',
	tagName: 'div',
	render: function() {
		this.$el.html(this.model.get('queue')[this.model.get('current')]);
		
		var busy = $('#busy');
		if(this.model.get('current') < this.model.get('queue').length - 1) {
			this.$el.append($('<div class="next"></div>'));
			busy.css('visibility', 'visible');
		} else {
			busy.css('visibility', 'hidden');
		}
		busy.off('click').on('click', function() {
			if(state.get('busy')) {
				game.get('message').next();
			}
		});
		
		return this;
	},
	initialize: function() {
		this.listenTo(this.model, 'change', this.render);
	}
});

/**
 * 下画面（シーン）
 * model: State
 */
var SceneView = Backbone.View.extend({
	id: 'scene',
	tagName: 'div',
	initialize: function() {
		this.listenTo(this.model, 'change:currentScene', this.sceneHasEntered);
	},
	render: function() {
		this.$el.empty();
		
		var scene = this.model.get('currentScene');
		this.$el.css('background-image', 'url("/download?blobkey=' + scene.get('background') + '")');

		scene.get('eventList').each(function(event) {
			var eventView = new EventView({
				model: event
			});
			this.$el.append(eventView.render().el);
		}, this);
		
		return this;
	},
	
	/**
	 * シーン開始時の処理
	 */
	sceneHasEntered: function() {
		var scene = this.model.get('currentScene');
		var code = scene.get('enter');
		if(code != '') {
			var event = new Event({
				code: JSON.parse(code)
			});
			event.execute();
		}
		
		this.render();
	}
});

/**
 * シーン内のイベント
 * model: Event
 */
var EventView = Backbone.View.extend({
	tagName: 'div',
	className: 'event',
	render: function() {
		if(this.model.get('visible') && this.model.get('removed') == false) {
			this.$el.show();
		} else {
			this.$el.hide();
			return this;
		}
		
		var position = this.model.get('position');
		this.$el.css('left', position[0]).css('top', position[1]);

		var size = this.model.get('size');
		this.$el.css('width', size[0]).css('height', size[1]);
		
		var blobkey = this.model.get('image');
		if(blobkey != '') {
			this.$el.css('background-image', 'url("/download?blobkey=' + blobkey + '")');
		}
		
		return this;
	},
	initialize: function() {
		this.listenTo(this.model, 'change', this.render);
	},
	events: {
		'click': 'eventHasClicked'
	},
	eventHasClicked: function() {
		if(state.get('busy')) {
			return;
		}
		this.model.execute();
		
		return false;
	}
});
