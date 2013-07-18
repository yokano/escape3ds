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
		this.on('change:name change:img change:hasFirst change:sort', this.itemHasChanged);
	},
	
	comparator: function(item) {
		return item.get('sort');
	},
	
	/**
	 * アイテムが選択された
	 */
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
	
	/**
	 * アイテムが追加された
	 */
	itemHasAdded: function(item) {
		Backbone.sync('create', item, {
			success: function(data) {
				item.set('id', data.itemKey);
			},
			error: function() {
				console.log('error');
			}
		});
	},
	
	/**
	 * アイテムが削除された
	 */
	itemHasRemoved: function(item) {
		Backbone.sync('delete', item, {
			error: function() {
				console.log('error');
			}
		});
	},
	
	/**
	 * アイテムの内容が更新された
	 */
	itemHasChanged: function(item) {
		Backbone.sync('update', item, {
			error: function() {
				console.log('error');
			}
		});
	}
});