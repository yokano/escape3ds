/**
 * シーンリストモデル
 */
var SceneList = Backbone.Collection.extend({
	model: Scene,
	selected: null,
	select: function(cid) {
		this.selected = cid;
	},
	removed: function(scene) {
		this.selected = null;
		Backbone.sync('delete', scene, {
			success: function() {
				console.log('success');
			},
			error: function() {
				console.log('error');
			}
		})
	},
	added: function(scene) {
		// シーンが追加されたらデータベースへ追加する
		Backbone.sync('create', scene, {
			success: function(data) {
				scene.set('id', data.sceneKey);
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