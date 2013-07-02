/**
 * アイテムモデル
 */
var Item = Backbone.Model.extend({
	defaults: {
		name: '',
		image: '',
		gameKey: ''
	}
});