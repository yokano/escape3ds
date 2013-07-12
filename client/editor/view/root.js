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
	mode: 'scene_editor',  // 表示中のエディタ
	render: function() {
		this.$el.empty();
		
		var headerView = new HeaderView({model: game});
		this.$el.append(headerView.render().el);
		
		switch(this.mode) {
		case 'scene_editor':
			var sceneEditorView = new SceneEditorView();
			this.$el.append(sceneEditorView.render().el);
			break;
		case 'item_editor':
			var itemEditorView = new ItemEditorView();
			this.$el.append(itemEditorView.render().el);
			break;
		}
		
		return this;
	},
	
	/**
	 * 表示モードを変更する
	 * @param {string} mode 表示モード名(scene_editor/item_editor)
	 */
	changeMode: function(mode) {
		if(this.mode != mode) {
			this.mode = mode;
			this.render();
		}
	}
});