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
	initialize: function(attr) {
		this.urlRoot = '/sync/scene/' + GAME_ID;
		this.on('change:name change:background', this.sceneHasChanged);

		if(attr.eventList == undefined) {
			this.set('eventList', new EventList());
		}
	}
});
