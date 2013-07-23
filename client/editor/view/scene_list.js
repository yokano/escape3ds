/**
 * シーンリストビュー
 * @class
 * @extends Backbone.View
 */
var SceneListView = Backbone.View.extend({
	tagName: 'ul',
	id: 'scene_list',
	initialize: function() {
		this.listenTo(this.collection, 'add', this.render);
		this.listenTo(this.collection, 'select', this.select);
	},
	render: function() {
		this.$el.empty();
		this.collection.each(function(scene) {
			var sceneView = new SceneListItemView({
				model: scene,
				parent: this
			});
			this.$el.append(sceneView.render().el);
		}, this);
		
		var self = this;
		this.$el.sortable({
			stop: function() {
				self.sortHasStopped();
			}
		});
		return this;
	},
	select: function(cid) {
		var scene = this.collection.get(cid);
		var index = this.collection.indexOf(scene)
		this.$el.children().removeClass('select');
		this.$el.children().eq(index).addClass('select');
	},
	sortHasStopped: function() {
		var self = this;
		this.$el.children().each(function(index) {
			var cid = $(this).attr('id');
			self.collection.get(cid).set('sort', index);
		});
		this.collection.sort();
	}
});
