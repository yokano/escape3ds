/**
 * アプリケーションのトップに表示されるビュー
 */
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

/**
 * ブロッグをドラッグする先の四角形
 * model: blockView のいずれか
 * ブロックをドラッグされたらイベントリストに追加する
 */
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
			accept: '.block',
			drop: function(event, ui) {
				blockList.remove(ui.draggable.attr('cid'), {'silent': true});
				view.eventList.add(new MethodBlock({type: ui.draggable.attr('type')}), {at: view.index});
			},
			hoverClass: 'connector-hover',
			over: function() {
				$('body').css('cursor', 'auto');
			},
			out: function() {
				$('body').css('cursor', 'url("/client/event_editor/trashbox.png"), auto');
			}
		});
		this.$el.html('ここへドラッグ');
		
		return this;
	}
});

/**
 * 画面左側の流れ図
 * collection: blockList
 */
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

/**
 * 画面右側のブロックを置いておくビュー
 */
var BlockListView = Backbone.View.extend({
	tagName: 'div',
	id: 'block_list_view',
	render: function() {
		var view = this;
		this.$el.append('<span class="label">ブロックリスト</span>');
		this.$el.append('<div class="block" type="message">メッセージ表示</div>');
		this.$el.append('<div class="block" type="changeScene">シーン移動</div>');
		this.$el.append('<div class="block" type="addItem">アイテム追加</div>');
		this.$el.append('<div class="block" type="removeItem">アイテム削除</div>');
		this.$el.append('<div class="block" type="hide">非表示</div>');
		this.$el.append('<div class="block" type="show">表示</div>');
		this.$el.append('<div class="block" type="remove">消滅</div>');
		this.$el.append('<div class="block" type="changeImage">画像変更</div>');
		this.$el.append('<div class="block" type="variable">変数操作</div>');
		
		this.$el.find('.block').draggable({
			helper: 'clone'
		});
		
		var closeButton = $('<button class="close">閉じる</button>');
		closeButton.on('click', function() {
			window.close();
		});
		this.$el.append(closeButton);
		
		return this;
	}
});

/**
 * 流れ図を構成するブロックのベースになるクラス
 * イベントの種類に合わせてそれぞれ継承先で template を準備すること
 */
var BlockView = Backbone.View.extend({
	tagName: 'div',
	render: function() {
		var view = this;
		var connectorView = new ConnectorView({
			eventList: blockList,
			index: blockList.indexOf(this.model) + 1
		});
		
		var block = $(this.template());
		block.attr('cid', view.model.cid);
		block.draggable({
			start: function(event, ui) {
				view.$el.find('.line').remove();
				connectorView.remove();
				block.removeClass('stack');
			},
			stop: function() {
				// 何もない場所にドラッグされた
				view.$el.fadeOut(function() {
					view.remove();
					blockList.remove(view.model);
					$('body').css('cursor', 'auto'); // jQuery UI が body の cursor を書き換えるため
				});
			},
			cursor: 'url("/client/event_editor/trashbox.png"), auto'
		});
		
		// ドラッグ開始直前にjQueryUIの座標のずれを修正する
		block.on('mousedown', function() {
			var difference = $(this).offset().left + $(this).width() / 2;
			$(this).draggable('option', 'cursorAt', {right: difference});
		});
		
		this.$el.append(block);
		this.$el.append('<div class="stack line"></div>');
		this.$el.append(connectorView.render().el);
		this.$el.append('<div class="stack line"></div>');
		
		return this;
	}
});

/**
 * シーンの変更ブロック
 */
var ChangeSceneBlockView = BlockView.extend({
	template: _.template($('#change_scene_template').html())
});

/**
 * アイテム追加ブロック
 */
var AddItemBlockView = BlockView.extend({
	template: _.template($('#add_item_template').html())
});

/**
 * アイテム削除ブロック
 */
var RemoveItemBlockView = BlockView.extend({
	template: _.template($('#remove_item_template').html())
});

/**
 * メッセージ表示ブロック
 */
var MessageBlockView = BlockView.extend({
	template: _.template($('#message_template').html())
});

/**
 * イベント非表示ブロック
 */
var HideBlockView = BlockView.extend({
	template: _.template($('#hide_template').html())
});

/**
 * イベント表示ブロック
 */
var ShowBlockView = BlockView.extend({
	template: _.template($('#show_template').html())
});

/**
 * イベント削除ブロック
 */
var RemoveBlockView = BlockView.extend({
	template: _.template($('#remove_template').html())
});

/**
 * 画像変更ブロック
 */
var ChangeImageBlockView = BlockView.extend({
	template: _.template($('#change_image_template').html())
});

/**
 * 変数操作ブロック
 */
var VariableBlockView = BlockView.extend({
	template: _.template($('#variable_template').html())
});

/**
 * ブロックの種類リスト
 */
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
