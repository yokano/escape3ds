/**
 * シーンモデル
 */
var Scene = Backbone.Model.extend({
	defaults: {
		name: '',
		gameKey: '',
		background: '/client/img/scene/black.png',
		enter: null,
		leave: null,
		eventList: null,
	},
	initialize: function() {
		this.set('eventList', new EventList());
	}
});