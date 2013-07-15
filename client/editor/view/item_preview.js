var ItemPreview = Backbone.View.extend({
	tagName: 'div',
	id: 'item_preview',
	template: _.template($('#item_preview_template').html()),
	initialize: function() {
		this.listenTo(this.model, 'change:hasFirst remove change:img', this.render);
	},
	render: function() {
		var firstItems = this.model.filter(function(item) {
			return item.get('hasFirst');
		});
		
		this.$el.html(this.template({items: firstItems}));
		return this;
	}
});