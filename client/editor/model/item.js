/**
 * アイテムモデル
 */
var Item = Backbone.Model.extend({
	defaults: {
		name: '',
		img: '',
		hasFirst: false,
		selected: false  // 現在選択されているかどうか
	}
});