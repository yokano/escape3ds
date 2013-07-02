/**
 * シーンビューの中に表示されるイベントや範囲
 * @class
 * @extends Backbone.View
 */
var EventAreaView = Backbone.View.extend({
	tagName: 'div',
	className: 'event_area',
	render: function() {
		return this;
	}
});