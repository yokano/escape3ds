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
			
		});
		this.$el.append(messageView.render().el);
		
		return this;
	}
});

/**
 * アイテムリスト
 */
var ItemListView = Backbone.View.extend({
	id: 'item_list',
	tagName: 'ul',
	render: function() {

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
	}
});

/**
 * 各アイテム
 */
var ItemView = Backbone.View.extend({
	className: 'item',
	tagName: 'li',
	render: function() {
		this.$el.append('<img src="' + this.model.get('img') + '">');
		return this;
	}
});

/**
 * メッセージ
 */
var MessageView = Backbone.View.extend({
	id: 'message',
	tagName: 'div',
	render: function() {
		return this;
	}
});

/**
 * 下画面（シーン）
 * model: State
 */
var SceneView = Backbone.View.extend({
	id: 'scene',
	tagName: 'div',
	render: function() {
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
	initialize: function() {
		this.listenTo(this.model, 'change:currentScene', this.render);
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
		if(this.model.get('visible')) {
			this.$el.show();
		} else {
			this.$el.hide();
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
		this.listenTo(this.model, 'change:visible', this.render);
	},
	events: {
		'click': 'eventHasClicked'
	},
	eventHasClicked: function() {
		var callback = new Function(this.model.get('code'));
		callback.call(this.model);
	}
});