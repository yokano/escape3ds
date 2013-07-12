/**
 * アイテム管理画面の中央下に表示されるアイテムの設定画面
 * @class
 */
var ItemInfoView = Backbone.View.extend({
	tagName: 'div',
	id: 'item_info',
	template: _.template($('#item_info_template').html()),
	render: function() {
		this.$el.html(this.template());
		return this;
	}
});