var ItemPreview = Backbone.View.extend({
	tagName: 'div',
	id: 'item_preview',
	template: _.template($('#item_preview_template').html()),
	render: function() {
		this.$el.html(this.template());
		return this;
	}
});