/**
 * シーンビューの中に表示されるイベントや範囲
 * @class
 * @extends Backbone.View
 */
var EventAreaView = Backbone.View.extend({
	tagName: 'div',
	className: 'event_area',
	initialize: function() {
		this.listenTo(this.model, 'change:selected', this.eventSelectedHasChanged);
	},
	render: function() {
		var model = this.model.toJSON();
		this.$el.css('left', model.position[0])
				.css('top', model.position[1])
				.css('width', model.size[0])
				.css('height', model.size[1]);
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
	}
});