/**
 * イベントモデル
 */
var Event = Backbone.Model.extend({
	defaults: {
		name: '',
		image: '',
		code: '',
		position: '',
		size: '',
		sceneKey: '',
		color: 'blue',
		selected: false
	},
	initialize: function() {
		console.log(this.attributes);
		this.urlRoot = '/sync/event/' + this.get('sceneId');
		this.on('change', this.eventHasChanged);
	},
	eventHasChanged: function() {
		Backbone.sync('update', this, {
			success: function() {
			},
			error: function() {
				console.log('error');
			}
		});
	}
});