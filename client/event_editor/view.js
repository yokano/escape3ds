var RootView = Backbone.View.extend({
	tagName: 'div',
	id: 'root',
	render: function() {
		var eventView = new EventView();
		this.$el.append(eventView.render().el);
		
		var blockListView = new BlockListView();
		this.$el.append(blockListView.render().el);
		
		return this;
	}
});

var EventView = Backbone.View.extend({
	tagName: 'div',
	id: 'event_view',
	render: function() {
		this.$el.append('<div class="label">イベント</div>');
		
		this.$el.append('<div class="stack circle" id="start">S</div>');
		this.$el.append('<div class="stack line"></div>');
		this.$el.append('<div class="stack connector">ここへドラッグ</div>');
		this.$el.append('<div class="stack line"></div>');
		this.$el.append('<div class="stack circle" id="end">E</div>');
		
		this.$el.find('.connector_area').droppable();
		
		return this;
	}
});

var BlockListView = Backbone.View.extend({
	tagName: 'div',
	id: 'block_list_view',
	render: function() {
		this.$el.append('<span class="label">ブロックリスト</span>');
		
		this.$el.append('<div class="block">シーンの移動</div>');
		this.$el.append('<div class="block">アイテムを追加</div>');
		this.$el.append('<div class="block">アイテムを削除</div>');
		this.$el.append('<div class="block">メッセージを表示</div>');
		this.$el.append('<div class="block">イベントを非表示にする</div>');
		this.$el.append('<div class="block">イベントを表示する</div>');
		this.$el.append('<div class="block">イベントの画像を変える</div>');
		this.$el.append('<div class="block">変数を操作する</div>');
		
		this.$el.find('.block').draggable({
			helper: 'clone',
			snap: '.connector',
			snapMode: 'inner'
		});
		
		return this;
	}
});