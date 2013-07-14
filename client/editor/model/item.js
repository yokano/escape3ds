/**
 * アイテムモデル
 */
var Item = Backbone.Model.extend({
	defaults: {
		name: '',
		img: '',
		selected: false  // 現在選択されているかどうか
	}
});