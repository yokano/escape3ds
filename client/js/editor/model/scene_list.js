/**
 * シーンリストモデル
 */
var SceneList = Backbone.Collection.extend({
	model: Scene,
	selected: null,
	select: function(cid) {
		this.selected = cid;
	},
	removed: function() {
		this.selected = null;
	},
	initialize: function() {
		this.on('select', this.select);
		this.on('remove', this.removed);
	},
	getSelected: function() {
		return this.get(this.selected);
	}
});