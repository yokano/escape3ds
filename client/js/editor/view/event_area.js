/**
 * シーンビューの中に表示されるイベントや範囲
 * @class
 * @extends Backbone.View
 */
var EventAreaView = Backbone.View.extend({
	tagName: 'div',
	className: 'event_area',
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
		this.$el.toggleClass('selected');
	}
});