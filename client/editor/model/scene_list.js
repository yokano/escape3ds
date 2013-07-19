/**
 * シーンリストモデル
 */
var SceneList = Backbone.Collection.extend({
	model: Scene,
	selected: null,
	comparator: function(scene) {
		return scene.get('sort');
	},
	removed: function(scene) {
		this.selected = null;
		Backbone.sync('delete', scene, {
			success: function() {
			},
			error: function() {
				console.log('error');
			}
		})
	},
	select: function(cid) {
		this.selected = cid;
	},
	added: function(scene) {
		// シーンが追加されたらデータベースへ追加する
		Backbone.sync('create', scene, {
			success: function(data) {
				scene.id = data.sceneKey;
			},
			error: function() {
				console.log('error');
			}
		});
	},
	initialize: function() {
		this.on('select', this.select);
		this.on('remove', this.removed);
		this.on('add', this.added);
	},
	getSelected: function() {
		return this.get(this.selected);
	}
});