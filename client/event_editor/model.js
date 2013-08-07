/**
 * ブロックの配列
 * MethodBlock と IfBlock を内包する
 */
var BlockList = Backbone.Collection.extend({
	url: '/sync/event/' + eventId,
	initialize: function(attr) {
		// DB から読み取った JSON を変換
		_.each(attr, function(block) {
			if(block.type == 'if') {
				var ifBlock = new IfBlock(block);
				this.add(ifBlock);
			} else {
				var methodBlock = new MethodBlock(block);
				this.add(methodBlock);
			}
		}, this);
		this.on('change add remove', this.update);
	},
	update: function() {
		$.ajax('/update_code', {
			method: 'POST',
			data: {
				id: eventId,
				code: JSON.stringify(this.toJSON())
			}
		});
	}
});

/**
 * 逐次命令ブロック
 * @property {String} type 命令の種類
 * @property {String] attr 命令の引数
 */
var MethodBlock = Backbone.Model.extend({
	defaults: {
		type: '',
		attr: ''
	}
});

/**
 * 条件分岐ブロック
 * @property {String} type if固定
 * @property {String} conditionType 条件の種類を表す。"hasItem" または "currentItem"。
 * @property {String} target 対象とするアイテムのキー。
 * @property {String} true 条件文がtrueの時に実行されるブロックリスト
 * @property {String} false 条件文がfalseの時に実行されるブロックリスト
 */
var IfBlock = Backbone.Model.extend({
	defaults: {
		type: 'if',
		conditionType: '',
		target: '',
		true: null,
		false: null
	},
	initialize: function(attr) {
//		this.set('true', new BlockList(attr.true));
//		this.set('false', new BlockList(attr.false));
		this.set('true', null);
		this.set('false', null);
		if(this.get('target') == '' && _.keys(itemList).length > 0) {
			this.set('target', _.keys(itemList)[0]);
		}
	}
});