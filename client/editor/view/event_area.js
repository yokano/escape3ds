/**
 * シーンビューの中に表示されるイベントや範囲
 * @class
 * @extends Backbone.View
 */
var EventAreaView = Backbone.View.extend({
	tagName: 'div',
	className: 'event_area',
	borderColor: {
		'red': 'yellow',
		'blue': 'red',
		'green': 'yellow',
		'yellow': 'red',
		'pink': 'red',
		'perple': 'yellow',
		'black': 'red',
		'white': 'red'
	},
	initialize: function() {
		this.listenTo(this.model, 'change:selected', this.eventSelectedHasChanged);
		this.listenTo(this.model, 'remove', this.eventHasRemoved);
		this.listenTo(this.model, 'change:color', this.colorHasChanged);
		this.listenTo(this.model, 'change:image', this.imageHasChanged);
	},
	render: function() {
		var model = this.model.toJSON();
		this.$el.css('left', model.position[0])
				.css('top', model.position[1])
				.css('width', model.size[0])
				.css('height', model.size[1])
				.css('background-color', model.color)
				.css('border-color', this.borderColor[model.color]);
		
		if(this.model.get('image') != '') {
			this.$el.css('background-image', 'url("/download?blobkey=' + this.model.get('image') + '")');
		}
		return this;
	},
	events: {
		'click': 'eventAreaHasClicked'
	},
	eventAreaHasClicked: function() {
		this.model.trigger('eventAreaHasSelected', this.model);
	},
	eventSelectedHasChanged: function() {
		if(this.model.get('selected')) {
			this.$el.addClass('selected');
		} else {
			this.$el.removeClass('selected');
		}
	},
	eventHasRemoved: function() {
		this.remove();
	},
	colorHasChanged: function() {
		var color = this.model.get('color');
		this.$el.css('background-color', color);
		this.$el.css('border-color', this.borderColor[color])
	},
	imageHasChanged: function() {
		this.render();
	}
});