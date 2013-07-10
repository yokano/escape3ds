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
	initialize: function() {
		this.urlRoot = '/sync/scene/' + game.id;
		this.set('eventList', new EventList());
	}
});
