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
	events: {
		'click #remove_event': 'removeButtonHasClicked'
	},
	initialize: function() {
		this.listenTo(game.get('sceneList'), 'select', this.sceneHasSelected);
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
		this.stopListening(this.collection);
		this.collection = game.get('sceneList').getSelected().get('eventList');
		this.listenTo(this.collection, 'eventAreaHasSelected', this.eventHasSelected);
		this.listenTo(this.collection, 'remove', this.eventHasRemoved);
	},
	eventHasSelected: function(selectedEvent) {
		if(selectedEvent == this.model) {
			return;
		}
		this.model = selectedEvent;
		this.render();
	},
	removeButtonHasClicked: function() {
		if(!window.confirm('イベントを削除しますか？')) {
			return;
		}
		this.model.trigger('removeButtonHasClicked', this.model);
	},
	eventHasRemoved: function() {
		this.model = null;
		this.render();
	}
});