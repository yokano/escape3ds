/**
 * アイテムリストを構成する <li> 要素
 * @class
 */
var ItemListItem = Backbone.View.extend({
	tagName: 'li',
	className: 'item_list_item',
	template: _.template($('#item_list_item_template').html()),
	initialize: function() {
		this.listenTo(this.model, 'change', this.render);
	},
	render: function() {
		this.$el.html(this.template(this.model.toJSON()));
		
		if(this.model.get('selected')) {
			this.$el.addClass('selected');
		} else {
			this.$el.removeClass('selected');
		}
		
		return this;
	},
	events: {
		'click': 'clicked'
	},
	clicked: function() {
		this.model.set('selected', true);
	}
});