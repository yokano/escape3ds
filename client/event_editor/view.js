var RootView = Backbone.View.extend({
	tagName: 'div',
	id: 'root',
	render: function() {
		var eventView = new EventView({
			collection: blockList
		});
		this.$el.append(eventView.render().el);
		
		var blockListView = new BlockListView();
		this.$el.append(blockListView.render().el);
		
		return this;
	}
});

// model: blockView のいずれか
// ブロックをドラッグされたらイベントリストに追加する
var ConnectorView = Backbone.View.extend({
	tagName: 'div',
	className: 'stack connector',
	initialize: function(attr) {
		this.eventList = attr.eventList; // 自分が属するイベントリスト
		this.index = (attr.index != undefined) ? attr.index : 0; // ドロップされた時に何番目に挿入するか
	},
	render: function() {
		var view = this;
		
		// ドロップされた
		this.$el.droppable({
			drop: function(event, ui) {
				view.eventList.add(new MethodBlock({type: ui.draggable.attr('type')}), {at: view.index});
			},
			hoverClass: 'connector-hover'
		});
		this.$el.html('ここへドラッグ');
		
		return this;
	}
});

// collection: blockList
var EventView = Backbone.View.extend({
	tagName: 'div',
	id: 'event_view',
	initialize: function() {
		this.listenTo(this.collection, 'add', this.render);
	},
	render: function() {
		var view = this;
		
		// コネクタ
		var connectorView = new ConnectorView({
			eventList: blockList,
			index: 0
		});
		
		// イベント内容のブロック
		var blocks = $('<div></div>');
		this.collection.each(function(block) {
			var blockView = new BlockViewClasses[block.get('type')]({
				model: block
			});
			blocks.append(blockView.render().el);
		});
		
		// 表示
		this.$el.empty();
		this.$el.append('<div class="label">イベント</div>');
		this.$el.append('<div class="stack circle" id="start">S</div>');
		this.$el.append('<div class="stack line"></div>');
		this.$el.append(connectorView.render().el);
		this.$el.append('<div class="stack line"></div>');
		this.$el.append(blocks);
		this.$el.append('<div class="stack circle" id="end">E</div>');
				
		return this;
	}
});

var BlockListView = Backbone.View.extend({
	tagName: 'div',
	id: 'block_list_view',
	render: function() {
		var view = this;
		this.$el.append('<span class="label">ブロックリスト</span>');
		this.$el.append('<div class="block" type="changeScene">シーン移動</div>');
		this.$el.append('<div class="block" type="addItem">アイテム追加</div>');
		this.$el.append('<div class="block" type="removeItem">アイテム削除</div>');
		this.$el.append('<div class="block" type="hide">非表示</div>');
		this.$el.append('<div class="block" type="show">表示</div>');
		this.$el.append('<div class="block" type="remove">消滅</div>');
		this.$el.append('<div class="block" type="changeImage">画像変更</div>');
		this.$el.append('<div class="block" type="variable">変数操作</div>');
		
		this.$el.find('.block').draggable({
			helper: 'clone',
			snap: '.connector',
			snapMode: 'inner'
		});
		
		return this;
	}
});

var ChangeSceneBlockView = Backbone.View.extend({
	tagName: 'div',
	attributes: {
		type: 'changeScene'
	},
	render: function() {
		var connectorView = new ConnectorView({
			eventList: blockList,
			index: blockList.indexOf(this.model) + 1
		});
		
		this.$el.append('<div class="stack block">シーン移動</div>');
		this.$el.append('<div class="stack line"></div>');
		this.$el.append(connectorView.render().el);
		this.$el.append('<div class="stack line"></div>');
		
		return this;
	}
});

var AddItemBlockView = Backbone.View.extend({
	tagName: 'div',
	attributes: {
		type: 'addItem'
	},
	render: function() {
		var connectorView = new ConnectorView({
			eventList: blockList,
			index: blockList.indexOf(this.model) + 1
		});
	
		this.$el.append('<div class="stack block">アイテムを追加</div>');
		this.$el.append('<div class="stack line"></div>');
		this.$el.append(connectorView.render().el);
		this.$el.append('<div class="stack line"></div>');
		
		return this;
	},
});

var RemoveItemBlockView = Backbone.View.extend({
	tagName: 'div',
	attributes: {
		type: 'removeItem'
	},
	render: function() {
		var connectorView = new ConnectorView({
			eventList: blockList,
			index: blockList.indexOf(this.model) + 1
		});
	
		this.$el.append('<div class="stack block">アイテムを削除</div>');
		this.$el.append('<div class="stack line"></div>');
		this.$el.append(connectorView.render().el);
		this.$el.append('<div class="stack line"></div>');
		
		return this;
	}
});

var MessageBlockView = Backbone.View.extend({
	tagName: 'div',
	attributes: {
		type: 'message'
	},
	render: function() {
		var connectorView = new ConnectorView({
			eventList: blockList,
			index: blockList.indexOf(this.model) + 1
		});
	
		this.$el.append('<div class="stack block">メッセージを表示</div>');
		this.$el.append('<div class="stack line"></div>');
		this.$el.append(connectorView.render().el);
		this.$el.append('<div class="stack line"></div>');
		
		return this;
	}
});

var HideBlockView = Backbone.View.extend({
	tagName: 'div',
	attributes: {
		type: 'hide'
	},
	render: function() {
		var connectorView = new ConnectorView({
			eventList: blockList,
			index: blockList.indexOf(this.model) + 1
		});
		
		this.$el.append('<div class="stack block">イベントを隠す</div>');
		this.$el.append('<div class="stack line"></div>');
		this.$el.append(connectorView.render().el);
		this.$el.append('<div class="stack line"></div>');
		
		return this;
	}
});

var ShowBlockView = Backbone.View.extend({
	tagName: 'div',
	attributes: {
		type: 'show'
	},
	render: function() {
		var connectorView = new ConnectorView({
			eventList: blockList,
			index: blockList.indexOf(this.model) + 1
		});
		
		this.$el.append('<div class="stack block">イベントを表示</div>');
		this.$el.append('<div class="stack line"></div>');
		this.$el.append(connectorView.render().el);
		this.$el.append('<div class="stack line"></div>');
		
		return this;
	}
});

var RemoveBlockView = Backbone.View.extend({
	tagName: 'div',
	attributes: {
		type: 'remove'
	},
	render: function() {
		var connectorView = new ConnectorView({
			eventList: blockList,
			index: blockList.indexOf(this.model) + 1
		});
	
		this.$el.append('<div class="stack block">イベントを削除</div>');
		this.$el.append('<div class="stack line"></div>');
		this.$el.append(connectorView.render().el);
		this.$el.append('<div class="stack line"></div>');
		
		return this;
	}
});

var ChangeImageBlockView = Backbone.View.extend({
	tagName: 'div',
	attributes: {
		type: 'changeImage'
	},
	render: function() {
		var connectorView = new ConnectorView({
			eventList: blockList,
			index: blockList.indexOf(this.model) + 1
		});
	
		this.$el.append('<div class="stack block">イベント画像を変更</div>');
		this.$el.append('<div class="stack line"></div>');
		this.$el.append(connectorView.render().el);
		this.$el.append('<div class="stack line"></div>');
		
		return this;
	}
});

var VariableBlockView = Backbone.View.extend({
	tagname: 'div',
	attributes: {
		type: 'variable'
	},
	render: function() {
		var connectorView = new ConnectorView({
			eventList: blockList,
			index: blockList.indexOf(this.model) + 1
		});
		
		this.$el.append('<div class="stack block">変数を操作</div>');
		this.$el.append('<div class="stack line"></div>');
		this.$el.append(connectorView.render().el);
		this.$el.append('<div class="stack line"></div>');
		
		return this;
	}
});

var BlockViewClasses = {
	'changeScene': ChangeSceneBlockView,
	'addItem': AddItemBlockView,
	'removeItem': RemoveItemBlockView,
	'message': MessageBlockView,
	'hide': HideBlockView,
	'show': ShowBlockView,
	'remove': RemoveBlockView,
	'changeImage': ChangeImageBlockView,
	'variable': VariableBlockView
};
