/**
 * ゲームモデル
 */
var Game = Backbone.Model.extend({
	urlRoot: '/sync/game',
	defaults: {
		name: '',
		description: '',
		thumbnail: '',
		userKey: '',
		firstScene: null,
		sceneList: null,
		itemList: null
	},
	initialize: function() {
		this.set('sceneList', new SceneList());
		this.set('itemList', new ItemList());
	}
});
