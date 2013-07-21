/**
 * ルートビュー
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
 */
var SceneView = Backbone.View.extend({
	id: 'scene',
	tagName: 'div',
	render: function() {
		var scene = game.sceneList[this.model.get('currentScene')];
		this.$el.css('background-image', 'url("/download?blobkey=' + scene.background + '")');
		return this;
	}
});