/**
 * ビューの定義
 *
 * rootView
 *   HeaderView
 *   SceneEditorView
 *     SceneListView
 *       SceneListItemView
 *     SceneView
 *       EventAreaView
 *     EventEditorView
 */
var RootView = Backbone.View.extend({
	tagName: 'div',
	id: 'root_view',
	render: function() {
		var headerView = new HeaderView({model: game});
		var sceneEditorView = new SceneEditorView()
		this.$el.append(
			headerView.render().el,
			sceneEditorView.render().el
		);
		return this;
	}
});