/**
 * アイテムリストを構成する <li> 要素
 * @class
 */
var ItemListItem = Backbone.View.extend({
	tagName: 'li',
	className: 'item_list_item',
	template: _.template($('#item_list_item_template').html()),
	render: function() {
		this.$el.html(this.template());
		return this;
	}
});