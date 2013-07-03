/**
 * ゲームモデル
 */
var Game = Backbone.Model.extend({
	defaults: {
		name: '',
		description: '',
		thumbnail: '',
		userKey: '',
		firstScene: null,
		sceneList: null
	},
	initialize: function() {
		this.set('sceneList', new SceneList());
	}
});