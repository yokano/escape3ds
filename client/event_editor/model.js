/**
 * ブロックの配列
 * MethodBlock と IfBlock を内包する
 */
var BlockList = Backbone.Collection.extend({
	url: '/sync/event/' + eventId,
	initialize: function(models, code, root) {
		// DB から読み取った JSON を変換
		_.each(code, function(block) {
			if(block.type == 'if') {
				var ifBlock = new IfBlock({
					conditionType: block.conditionType,
					target: block.target,
					yes: new BlockList(null, block.yes),
					no: new BlockList(null, block.no)
				});
				this.add(ifBlock);
			} else {
				var methodBlock = new MethodBlock(block);
				this.add(methodBlock);
			}
		}, this);
		
		// 大元になるブロックリストだけを保存
		if(root) {
			this.on('change add remove', this.update);
		} else {
			this.on('change add remove', function() {
				blockList.update();
			});
		}
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
		target: ''
	},
	initialize: function() {
		if(this.get('target') == '' && _.keys(itemList).length > 0) {
			this.set('target', _.keys(itemList)[0]);
		}
	}
});