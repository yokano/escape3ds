/**
 * シーンモデル
 */
var Scene = Backbone.Model.extend({
	defaults: {
		name: '',
		background: '/client/editor/img/black.png',
		enter: null,
		leave: null,
		eventList: null,
	},
	sceneHasChanged: function() {
		Backbone.sync('update', this, {
			success: function() {
				console.log('success');
			},
			error: function() {
				console.log('error');
			}
		});
	},
	initialize: function() {
		this.urlRoot = '/sync/scene/' + game.id;
		this.set('eventList', new EventList());
		this.on('change', this.sceneHasChanged);
	}
});
