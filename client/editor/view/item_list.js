/**
 * アイテム一覧ビュー
 * アイテムエディタの画面左側に表示される
 * @class
 */
var ItemListView = Backbone.View.extend({
	tagName: 'ul',
	id: 'item_list',
	template: _.template($('#item_list_template').html()),
	initialize: function() {
		this.model = game.get('itemList');
		this.listenTo(this.model, 'change', function() {
			console.log('changed');
		});
	},
	render: function() {
		this.$el.html(this.template());
		
		var self = this;
		this.model.each(function(item) {
			var li = new ItemListItem({model: item});
			self.$el.append(li.render().el);
		});
		
		return this;
	}
});