/**
 * イベント編集ビュー
 * 画面右側
 * @class
 * @extends Backbone.View
 */
var EventEditorView = Backbone.View.extend({
	tagName: 'div',
	id: 'event_editor',
	template: _.template($('#event_view_template').html()),
	render: function() {
		this.$el.html(this.template());
		return this;
	}
});