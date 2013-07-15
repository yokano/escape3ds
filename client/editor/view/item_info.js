/**
 * アイテム管理画面の中央下に表示されるアイテムの設定画面
 * @class
 */
var ItemInfoView = Backbone.View.extend({
	tagName: 'div',
	id: 'item_info',
	template: _.template($('#item_info_template').html()),
	initialize: function() {
		this.model = null;
		this.listenTo(game.get('itemList'), 'change:selected', function(item) {
			if(item.get('selected')) {
				this.model = item;
				this.render();
			}
		});
	},
	events: {
		'click #delete_item': 'deleteItemButtonHasClicked',
		'click #has_first': 'hasFirstCheckboxHasClicked',
		'change #item_name': 'itemNameHasChanged'
	},
	render: function() {
		if(this.model == null) {
			this.$el.hide();
			return this;
		}
		
		this.$el.html(this.template(this.model.toJSON()));
		this.$el.find('.item_img').css('background-image', 'url("' + this.model.get('img') + '")');
		this.$el.show();
		
		return this;
	},
	deleteItemButtonHasClicked: function() {
		game.get('itemList').remove(this.model);
	},
	hasFirstCheckboxHasClicked: function() {
		this.model.set('hasFirst', !this.model.get('hasFirst'));
		if(this.model.get('hasFirst')) {
			this.$el.find('has_first').attr('checked', 'true');
		} else {
			this.$el.find('has_first').attr('checked', 'false');
		}
	},
	itemNameHasChanged: function(event) {
		this.model.set('name', $(event.target).val());
	}
});