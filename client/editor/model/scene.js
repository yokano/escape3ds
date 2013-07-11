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
		console.log('シーンが更新されました');
		Backbone.sync('update', this, {
			success: function() {
				console.log('success');
			},
			error: function() {
				console.log('error');
			}
		});
	},
	backgroundHasChanged: function() {
		console.log('background changed');
	},
	initialize: function() {
		this.urlRoot = '/sync/scene/' + game.id;
		this.set('eventList', new EventList());
		this.on('change', this.sceneHasChanged);
		this.on('change:background', this.backgroundHasChanged);
	}
});
