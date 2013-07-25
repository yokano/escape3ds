/**
 * シーンエディタビュー
 * シーンの編集を行うビュー
 * ヘッダの下に表示される
 * 左にシーンリスト、中央にシーン、右にイベントエディタを表示
 * @class
 * @extends Backbone.View
 */
var SceneEditorView = Backbone.View.extend({
	tagname: 'section',
	id: 'scene_editor',
	render: function() {
		var sceneListView = new SceneListView({collection: game.get('sceneList')});
		var sceneView = new SceneView();
		var eventEditorView = new EventEditorView();
		
		$('<button id="add_scene">シーンを追加</button>').appendTo(this.$el);
		this.$el.append(
			sceneListView.render().el,
			sceneView.render().el,
			eventEditorView.render().el
		);
		return this;
	},
	events: {
		'click #add_scene': 'addScene'
	},
	addScene: function() {
		game.get('sceneList').add({
			name: '新しいシーン',
			sort: game.get('sceneList').length
		});
	}
});
