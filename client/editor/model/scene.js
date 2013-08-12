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
	initialize: function(attr) {
		this.urlRoot = '/sync/scene/' + GAME_ID;
		this.on('change:name change:background', this.sceneHasChanged);
		this.on('remove', this.sceneHasRemoved);

		if(attr.eventList == undefined) {
			this.set('eventList', new EventList());
		}
	},
	sceneHasChanged: function() {
		if(game.get('firstScene') == this.id) {
			game.set('thumbnail', this.get('background'));
		}
		
		Backbone.sync('update', this, {
			success: function() {
			},
			error: function() {
				console.log('error');
			}
		});
	},
	sceneHasRemoved: function() {
		if(game.get('firstScene') == this.id) {
			game.set('firstScene', null)
		}
	}
});
