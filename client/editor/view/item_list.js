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
		this.listenTo(this.model, 'add remove', this.render);
		
		var self = this;
		this.$el.sortable({
			stop: function() {
				self.sortHasStopped();
			}
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
	},
	
	/**
	 * アイテムリストがドラッグで並び替えられた
	 */
	sortHasStopped: function() {
		var self = this;
		this.$el.children().each(function(index) {
			var cid = $(this).attr('id');
			self.model.get(cid).set('sort', index);
		});
	}
});