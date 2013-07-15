/**
 * アイテムリスト
 */
var ItemList = Backbone.Collection.extend({
	model: Item,
	initialize: function() {
		this.urlRoot = '/sync/item/' + GAME_ID;
		this.on('change:selected', this.itemHasSelected);
		this.on('add', this.itemHasAdded);
		this.on('remove', this.itemHasRemoved);
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
	},
	itemHasAdded: function(item) {
		Backbone.sync('create', item, {
			success: function() {
				console.log('success');
			},
			error: function() {
				console.log('error');
			}
		});
	},
	itemHasRemoved: function(item) {
		Backbone.sync('delete', item, {
			success: function() {
				console.log('success');
			},
			error: function() {
				console.log('error');
			}
		});
	}
});