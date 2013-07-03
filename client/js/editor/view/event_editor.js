/**
 * イベント編集ビュー
 * 画面右側
 * @class
 * @extends Backbone.View
 */
var EventEditorView = Backbone.View.extend({
	tagName: 'div',
	id: 'event_editor',
	eventList: null,
	template: _.template($('#event_view_template').html()),
	initialize: function() {
		this.listenTo(sceneList, 'select', this.sceneHasSelected);
	},
	render: function() {
		if(this.model == null) {
			this.$el.empty();
			return this;
		}
		this.$el.html(this.template(this.model.toJSON()));

		return this;
	},
	sceneHasSelected: function() {
		if(this.eventList != null) {
			this.stopListening(this.eventList, 'eventAreaHasSelected');
			this.model = null;
			this.render();
		}
		this.eventList = sceneList.getSelected().get('eventList');
		this.listenTo(this.eventList, 'eventAreaHasSelected', this.eventHasSelected);
	},
	eventHasSelected: function(selectedEvent) {
		this.model = selectedEvent;
		this.render();
	}
});