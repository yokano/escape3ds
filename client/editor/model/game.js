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
		this.on('change:firstScene change:thumbnail', this.update);
	},
	update: function() {
		if(this.get('firstScene') == '') {
			this.set('thumbnail', '', {silent: true});
		}
	
		Backbone.sync('update', this, {
			success: function() {
			},
			error: function() {
				console.log('error');
			}
		});
	}
});
