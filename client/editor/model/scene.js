/**
 * シーンモデル
 */
var Scene = Backbone.Model.extend({
	defaults: {
		name: '',
		background: '',  // blobkey又は画像のURL
		enter: null,
		leave: null,
		eventList: null
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
	initialize: function(attr) {
		this.urlRoot = '/sync/scene/' + GAME_ID;
		this.on('change:name', this.sceneHasChanged);
		this.on('change:background', this.backgroundHasChanged);

		if(attr.eventList == undefined) {
			this.set('eventList', new EventList());
		}
	}
});
