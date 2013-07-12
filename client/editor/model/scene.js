/**
 * シーンモデル
 */
var Scene = Backbone.Model.extend({
	defaults: {
		name: '',
		background: '',  // blobkey又は画像のURL
		enter: null,
		leave: null,
		eventList: null,
	},
	sceneHasChanged: function() {
		Backbone.sync('update', this, {
			success: function() {
			},
			error: function() {
				console.log('error');
			}
		});
	},
	backgroundHasChanged: function() {
	},
	initialize: function() {
		this.urlRoot = '/sync/scene/' + game.id;
		this.set('eventList', new EventList());
		this.on('change', this.sceneHasChanged);
		this.on('change:background', this.backgroundHasChanged);
	}
});
