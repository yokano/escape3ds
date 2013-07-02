/**
 * ゲームモデル
 */
var Game = Backbone.Model.extend({
	defaults: {
		name: '',
		description: '',
		thumbnail: '',
		userKey: '',
		firstScene: null
	}
});