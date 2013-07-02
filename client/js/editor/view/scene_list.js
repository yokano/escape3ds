/**
 * シーンリストビュー
 * @class
 * @extends Backbone.View
 */
var SceneListView = Backbone.View.extend({
	tagName: 'ul',
	id: 'scene_list',
	initialize: function() {
		this.listenTo(this.collection, 'add remove', this.render);
		this.listenTo(this.collection, 'select', this.select);
		this.listenTo(game, 'change', this.gameHasChanged);
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
		return this;
	},
	select: function(cid) {
		var scene = this.collection.get(cid);
		var index = this.collection.indexOf(scene)
		this.$el.children().removeClass('select');
		this.$el.children().eq(index).addClass('select');
	},
	gameHasChanged: function() {
		this.$el.find('.is_first_scene').hide();
		var index = this.collection.indexOf(game.get('firstScene'));
		if(index == -1) {
			return;
		}
		this.$el.children().eq(index).find('.is_first_scene').show();
	}
});
