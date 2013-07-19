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
		'click #remove_event': 'removeButtonHasClicked',
		'change #color': 'colorHasChanged',
		'change #change_event_img': 'eventImageFileHasChanged',
		'click #remove_img': 'removeImageButtonHasClicked',
		'change #event_code': 'eventCodeHasChanged',
		'change .event_name': 'eventNameHasChanged'
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
		this.$el.find('#color').val(this.model.get('color'));
		
		return this;
	},
	sceneHasSelected: function() {
		if(this.collection != undefined) {
			this.stopListening(this.collection);
		}
		this.collection = game.get('sceneList').getSelected().get('eventList');
		this.listenTo(this.collection, 'eventAreaHasSelected', this.eventHasSelected);
		this.listenTo(this.collection, 'remove', this.eventHasRemoved);
		this.$el.empty();
	},
	eventHasSelected: function(selectedEvent) {
		this.model = selectedEvent;
		this.render();
	},
	removeButtonHasClicked: function() {
		if(!window.confirm('イベントを削除しますか？')) {
			return;
		}
		
		game.get('sceneList').getSelected().get('eventList').remove(this.model);
	},
	eventHasRemoved: function() {
		this.model = null;
		this.render();
	},
	colorHasChanged: function() {
		var color = this.$el.find('#color').val();
		this.model.set('color', color);
	},
	eventImageFileHasChanged: function() {
		var file = $('#change_event_img').get(0).files[0];
		getFileURL(file, this, function(url) {
			this.model.set('image', url);
		});
	},
	removeImageButtonHasClicked: function() {
		this.$el.find('#change_event_img').val('');
		this.model.set('image', '');
	},
	eventCodeHasChanged: function(event) {
		var code = $(event.target).val();
		this.model.set('code', code);
	},
	eventNameHasChanged: function(event) {
		var name = $(event.target).val();
		this.model.set('name', name);
	}
});