/**
 * アイテムリスト
 */
var ItemList = Backbone.Collection.extend({
	model: Item,
	initialize: function() {
		this.on('change:selected', this.itemHasSelected);
	},
	itemHasSelected: function(selected) {
		if(selected.get('selected') == false) {
			return;
		}
		
		this.each(function(item) {
			if(item.cid != selected.cid) {
				item.set('selected', false);
			}
		});
	}
});