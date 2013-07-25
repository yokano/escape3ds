/**
 * アイテムエディタビュー
 * アイテム管理ボタンを押した時に表示される
 * @class
 */
var ItemEditorView = Backbone.View.extend({
	tagName: 'div',
	id: 'item_editor',
	template: _.template($('#item_editor_template').html()),
	render: function() {
		this.$el.html(this.template());
		
		var itemListView = new ItemListView();
		this.$el.append(itemListView.render().el);
		
		var itemPreview = new ItemPreview({model: game.get('itemList')});
		this.$el.append(itemPreview.render().el);
		
		var itemInfoView = new ItemInfoView();
		this.$el.append(itemInfoView.render().el);
		
		return this;
	},
	events: {
		'click #add_item': 'addItemButtonHasClicked'
	},
	
	/**
	 * アイテムを追加ボタンが押された
	 */
	addItemButtonHasClicked: function() {
		game.get('itemList').add({
			name: '新しいアイテム',
			sort: game.get('itemList').models.length
		});
	}
})