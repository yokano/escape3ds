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
		this.urlRoot = '/sync/event/' + this.get('sceneId');
		this.on('change:name change:image change:position change:size change:color', this.eventHasChanged);
	},
	eventHasChanged: function() {
		console.log('update');
		console.log(this.changedAttributes());
		Backbone.sync('update', this, {
			success: function() {
			},
			error: function() {
				console.log('error');
			}
		});
	}
});