/**
 * イベントモデル
 */
var Event = Backbone.Model.extend({
	defaults: {
		name: '',
		image: '',
		code: '',
		position: '',
		size: '',
		sceneKey: '',
		color: 'blue',
		selected: false
	}
});