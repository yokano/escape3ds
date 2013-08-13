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
	},
	render: function() {
		var view = this;
		
		// ドロップされた
		this.$el.droppable({
			accept: '.block',
			drop: function(event, ui) {
				// このコネクタを所有しているブロック
				var model = ui.draggable.data('model');
				
				// ドラッグされたモデルが所属しているブロックリスト
				var collection = ui.draggable.data('collection');
				
				// 移動元のブロックを削除する
				if(collection != undefined) {
					collection.remove(model, {silent: true});
				}
				
				// このコネクタにドラッグされた時にどこへ挿入するか
				if(view.model == null) {
					view.index = 0;
				} else {
					view.index = view.eventList.indexOf(view.model) + 1;
				}
				
				// 画面右のブロックリストからドラッグされた
				if(model == undefined) {
					var type = ui.draggable.attr('type');
					if(type == 'if') {
						var ifBlock = new IfBlock({
							type: 'if',
							conditionType: ui.draggable.find('.conditionType').val(),
							target: ui.draggable.find('.target').val(),
							yes: new BlockList(),
							no: new BlockList()
						});
						view.eventList.add(ifBlock, {at: view.index});
					} else {
						var methodBlock = new MethodBlock({
							type: ui.draggable.attr('type')
						});
						view.eventList.add(methodBlock, {at: view.index});
					}
				} else {
					// 配置済みのブロックがドラッグされた
					view.eventList.add(model, {at: view.index});
				}
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
		this.listenTo(this.collection, 'add remove reset', this.render);
	},
	render: function() {
		var view = this;
		
		// コネクタ
		var connectorView = new ConnectorView({
			eventList: blockList,
			model: null
		});
		
		// イベント内容のブロック
		var blocks = $('<div></div>');
		this.collection.each(function(block) {
			var blockView = new BlockViewClasses[block.get('type')]({
				model: block,
				blockList: blockList
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
		this.$el.append('<div class="block method" type="message">メッセージ表示</div>');
		this.$el.append('<div class="block method" type="changeScene">シーン移動</div>');
		this.$el.append('<div class="block method" type="addItem">アイテム追加</div>');
		this.$el.append('<div class="block method" type="removeItem">アイテム削除</div>');
		this.$el.append('<div class="block method" type="hide">非表示</div>');
		this.$el.append('<div class="block method" type="show">表示</div>');
		this.$el.append('<div class="block method" type="remove">消滅</div>');
		this.$el.append('<div class="block method" type="changeImage">画像変更</div>');
		this.$el.append('<div class="block if" type="if">条件分岐</div>');
		
		this.$el.find('.block').draggable({
			helper: 'clone',
			zIndex: 3
		});
		
		var closeButton = $('<button class="close">閉じる</button>');
		closeButton.on('click', function() {
			window.close();
		});
		this.$el.append(closeButton);
		
		var clearButton = $('<button class="clear">クリア</button>');
		clearButton.on('click', function() {
			if(window.confirm('イベント内容をすべて削除しますか？')) {
				blockList.reset([]);
			}
		});
		this.$el.append(clearButton);
		
		return this;
	}
});

/**
 * 流れ図を構成するブロックのベースになるクラス
 * イベントの種類に合わせてそれぞれ継承先で template を準備すること
 * @property {BlockList} blockList このブロックが属するブロックリストへの参照
 * 条件分岐によって分岐された命令はそれぞれのブロックリストに属する
 */
var MethodBlockView = Backbone.View.extend({
	tagName: 'div',
	initialize: function(options) {
		this.blockList = options.blockList;
	},
	render: function() {
		var view = this;
		var connectorView = new ConnectorView({
			eventList: this.blockList,
			model: this.model
		});
		
		var block = $(this.template());
		block.draggable({
			start: function(event, ui) {
				$(this).data('model', view.model);
				$(this).data('collection', view.blockList);

				view.$el.find('.line').remove();
				connectorView.remove();
				block.removeClass('stack');
			},
			stop: function() {
				// 何もない場所にドラッグされた
				view.$el.fadeOut(function() {
					view.remove();
					view.options.blockList.remove(view.model);
					$('body').css('cursor', 'auto'); // jQuery UI が body の cursor を書き換えるため
				});
			},
			cursor: 'url("/client/event_editor/trashbox.png"), auto',
			zIndex: 3
		});
		
		// ドラッグ開始直前にjQueryUIの座標のずれを修正する
		block.on('mousedown', function() {
			var marginLeft = parseInt($(this).css('margin-left').slice(0, -2));
			var difference = marginLeft + $(this).width() / 2;
			$(this).draggable('option', 'cursorAt', {right: difference});
		});
		
		// シーンの選択肢を持っていたら表示
		var scene = block.find('.scene');
		if(scene.length > 0) {
			_.each(sceneList, function(val, key) {
				var option = $('<option></option>');
				option.html(val);
				option.val(key);
				if(key == view.model.get('attr')) {
					option.attr('selected', '');
				}
				scene.append(option);
			});
		}
		
		// アイテムの選択肢を持っていたら表示
		var item = block.find('.item');
		if(item.length > 0) {
			_.each(itemList, function(val, key) {
				var option = $('<option></option>');
				option.html(val);
				option.val(key);
				if(key == view.model.get('attr')) {
					option.attr('selected', '');
				}
				item.append(option);
			});
		}
		
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
var ChangeSceneBlockView = MethodBlockView.extend({
	template: _.template($('#change_scene_template').html()),
	constructor: function(options) {
		MethodBlockView.call(this, options);
		if(this.model.get('attr') == '') {
			this.model.set('attr', _.keys(sceneList)[0]);
		}
	},
	events: {
		'change .scene': 'sceneHasSelected'
	},
	sceneHasSelected: function() {
		var scene = this.$el.find('.scene').val();
		this.model.set('attr', scene);
	}
});

/**
 * アイテム追加ブロック
 */
var AddItemBlockView = MethodBlockView.extend({
	template: _.template($('#add_item_template').html()),
	constructor: function(options) {
		MethodBlockView.call(this, options);
		if(this.model.get('attr') == '') {
			this.model.set('attr', _.keys(itemList)[0]);
		}
	},
	events: {
		'change .item': 'itemHasSelected'
	},
	itemHasSelected: function() {
		var item = this.$el.find('.item').val();
		this.model.set('attr', item);
	}
});

/**
 * アイテム削除ブロック
 */
var RemoveItemBlockView = MethodBlockView.extend({
	template: _.template($('#remove_item_template').html()),
	constructor: function(options) {
		MethodBlockView.call(this, options);
		if(this.model.get('attr') == '') {
			this.model.set('attr', _.keys(itemList)[0]);
		}
	},
	events: {
		'change .item': 'itemHasChanged'
	},
	itemHasChanged: function() {
		var item = this.$el.find('.item').val();
		this.model.set('attr', item);
	}
});

/**
 * メッセージ表示ブロック
 */
var MessageBlockView = MethodBlockView.extend({
	template: _.template($('#message_template').html())
});

/**
 * イベント非表示ブロック
 */
var HideBlockView = MethodBlockView.extend({
	template: _.template($('#hide_template').html())
});

/**
 * イベント表示ブロック
 */
var ShowBlockView = MethodBlockView.extend({
	template: _.template($('#show_template').html())
});

/**
 * イベント削除ブロック
 */
var RemoveBlockView = MethodBlockView.extend({
	template: _.template($('#remove_template').html())
});

/**
 * 画像変更ブロック
 */
var ChangeImageBlockView = MethodBlockView.extend({
	template: _.template($('#change_image_template').html()),
	constructor: function(options) {
		MethodBlockView.call(this, options);
	},
	events: {
		'change .img': 'imageHasChanged'
	},
	imageHasChanged: function() {
		var view = this;

		// 古い画像を削除
		if(this.model.get('attr') != '') {
			$.ajax('/delete_blog', {
				method: 'POST',
				async: false,
				data: {
					blobkey: view.model.get('attr')
				},
				error: function() {
					console.log('画像の削除に失敗');
				}
			});
		}
		
		var file = this.$el.find('.img').get(0).files[0];
		var formdata = new FormData();
		formdata.append('file', file);
		
		// 新しい画像をアップロード
		$.ajax(geturl(), {
			method: 'POST',
			contentType: false,
			processData: false,
			data: formdata,
			dataType: 'json',
			error: function() {
				console.log('ファイルのアップロードに失敗しました');
			},
			success: function(data) {
				view.model.set('attr', data.blobkey);
			}
		});
	}
});

/**
 * 条件分岐ブロック
 */
var IfBlockView = Backbone.View.extend({
	tagName: 'div',
	className: 'stack if_container',
	template: _.template($('#if_template').html()),
	initialize: function(options) {
		this.blockList = options.blockList;
		this.listenTo(this.model.get('yes'), 'add remove', this.render);
		this.listenTo(this.model.get('no'), 'add remove', this.render);
	},
	events: {
		'change .target': 'targetHasSelected',
		'change .conditionType': 'conditionTypeHasSelected'
	},
	targetHasSelected: function(event) {
		this.model.set('target', $(event.target).val());
	},
	conditionTypeHasSelected: function(event) {
		this.model.set('conditionType', $(event.target).val());
	},
	render: function() {
		this.$el.empty();
		
		var view = this;
		var ifBlock = $('<div type="if" class="stack block if"></div>');
		ifBlock.html(this.template());
		ifBlock.find('.target').val(this.model.get('target'));
		
		// ドラッグ
		ifBlock.draggable({
			start: function() {
				$(this).data('model', view.model);
				$(this).data('collection', view.blockList);
				
				view.$el.find('.if_body').remove();
				view.$el.find('.if_footer').remove();
				view.$el.find('.left').remove();
				view.$el.find('.right').remove();
				view.$el.find('.start_if_line').remove();
				
				$(this).removeClass('stack');
			},
			stop: function() {
				view.$el.remove();
				view.blockList.remove(view.model);
			},
			cursor: 'url("/client/event_editor/trashbox.png"), auto',
			zIndex: 3
		});
		
		// ドラッグ時のズレを修正
		ifBlock.on('mousedown', function() {
			var marginLeft = $(this).css('margin-left').slice(0, -2);
			marginLeft = parseInt(marginLeft) + $(this).width() / 2;
			ifBlock.draggable('option', 'cursorAt', {right: marginLeft});
		});
		
		// アイテムリストを選択肢に追加
		_.each(itemList, function(item, key) {
			var option = $('<option value="' + key + '">' + item + '</option>');
			ifBlock.find('.target').append(option);
		}, this);
		ifBlock.find('.target [value="' + this.model.get('target') + '"]').attr('selected', '');
		
		// 条件の選択
		ifBlock.find('.conditionType [value="' + this.model.get('conditionType') + '"]').attr('selected', '');
		
		// 分岐開始
		var header = $('<div class="stack if_header"></div>');
		header.append(ifBlock);
		header.append('<div class="stack start_if_line"></div>');
		var left = $('<div class="if_line_container left"></div>');
		left.append('<div class="stack line"></div>');
		left.append('<div class="yes">YES</div>');
		header.append(left);
		var right = $('<div class="if_line_container right"></div>');
		right.append('<div class="stack line"></div>');
		right.append('<div class="no">NO</div>');
		header.append(right);
		this.$el.append(header);
		
		// 処理内容
		var body = $('<div class="stack if_body"></div>');
		
		// 左側の処理
		var left = $('<div class="if_container_left"></div>');
		var connectorView = new ConnectorView({
			eventList: this.model.get('yes'),
			model: null
		});
		left.append(connectorView.render().el);
		left.append('<div class="stack line"></div>');
		this.model.get('yes').each(function(block) {
			var blockView = new BlockViewClasses[block.get('type')]({
				model: block,
				blockList: this.model.get('yes')
			});
			left.append(blockView.render().el);
		}, this);
		left.append('<div class="stack line"></div>');
		body.append(left);
		
		// 右側の処理
		var right = $('<div class="if_container_right"></div>');
		var connectorView = new ConnectorView({
			eventList: this.model.get('no'),
			model: null
		});
		right.append(connectorView.render().el);
		right.append('<div class="stack line"></div>');
		this.model.get('no').each(function(block) {
			var blockView = new BlockViewClasses[block.get('type')]({
				model: block,
				blockList: this.model.get('no')
			});
			right.append(blockView.render().el);
		}, this);
		right.append('<div class="stack line"></div>');
		body.append(right);
		
		this.$el.append(body);
		
		// 分岐終了
		var footer = $('<div class="stack if_footer"></div>');
		footer.append('<div class="stack end_if_line"></div>');
		footer.append('<div class="if_line_container left"><div class="stack line"></div></div>');
		footer.append('<div class="if_line_container right"><div class="stack line"></div></div>');
		footer.append('<div class="stack line"></div>');
		var connectorView = new ConnectorView({
			eventList: this.blockList,
			model: this.model
		});
		footer.append(connectorView.render().el);
		footer.append('<div class="stack line"></div>');
		this.$el.append(footer);
		
		return this;
	}
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
	'if': IfBlockView
};
