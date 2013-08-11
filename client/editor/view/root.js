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
	jcropAPI: null,
	initialize: function() {
		$(window).on('keydown', this.keyHasDown);
	},
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
	},
	
	/**
	 * キーボードが押された
	 */
	 keyHasDown: function(e) {
		if((e.keyCode == 8 || e.keyCode == 46) && game.get('sceneList').getSelected().get('eventList').getSelected() != undefined) {
			$('#remove_event').click();
		}
	 }
});