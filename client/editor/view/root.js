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
		var itemEditorView = new ItemEditorView();
		this.$el.append(
			headerView.render().el,
//			sceneEditorView.render().el
			itemEditorView.render().el
		);
		return this;
	}
});