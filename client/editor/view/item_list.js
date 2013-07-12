/**
 * アイテム一覧ビュー
 * アイテムエディタの画面左側に表示される
 * @class
 */
var ItemListView = Backbone.View.extend({
	tagName: 'ul',
	id: 'item_list',
	template: _.template($('#item_list_template').html()),
	render: function() {
		this.$el.html(this.template());
		
		var li = new ItemListItem()
		this.$el.append(li.render().el);
		
		return this;
	}
});