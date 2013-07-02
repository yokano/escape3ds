/**
 * シーンモデル
 */
var Scene = Backbone.Model.extend({
	defaults: {
		name: '',
		gameKey: '',
		background: '/client/img/scene/black.png',
		enter: '',
		leave: ''
	}
});