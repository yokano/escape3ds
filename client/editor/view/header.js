/**
 * ヘッダービュー
 * 画面上に表示されるヘッダ部分のビュー
 * 表示するモデル : game
 * @class
 * @extends Backbone.View
 */
var HeaderView = Backbone.View.extend({
	tagName: 'header',
	template: _.template($('#header_view_template').html()),
	render: function() {
		this.$el.html(this.template(this.model.toJSON()));
		return this;
	},
	events: {
		'click #scene_mode' : 'changeToSceneMode',
		'click #item_mode': 'changeToItemMode',
		'click #back' : 'backToGameList',
		'click #save' : 'save'
	},
	changeToSceneMode: function() {
		rootView.changeMode('scene_editor');
	},
	changeToItemMode: function() {
		rootView.changeMode('item_editor');
	},
	backToGameList: function() {
		console.log('ゲームリストへ戻る');
	},
	save: function() {
		// ゲーム更新
		Backbone.sync('update', game, {
			success: function() {
			},
			error: function() {
				console.log('error');
			}
		});
	}
});