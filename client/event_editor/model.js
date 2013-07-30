var BlockList = Backbone.Collection.extend({
	url: '/sync/event/' + eventId,
	initialize: function() {
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

var MethodBlock = Backbone.Model.extend({
	defaults: {
		type: '',
		attr: ''
	}
});

/**
 * 条件分岐ブロック
 * @property {String} type if固定
 * @property {String} condition 条件文
 * @property {String} true 条件文がtrueの時に実行されるブロックリスト
 * @property {String} false 条件文がfalseの時に実行されるブロックリスト
 */
var IfBlock = Backbone.Model.extend({
	defaults: {
		type: 'if',
		condition: 'true',
		true: null,
		false: null
	}
});