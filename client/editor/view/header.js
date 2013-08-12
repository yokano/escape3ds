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
		
		if(rootView.mode == 'scene_editor') {
			this.$el.find('#scene_mode')
		}
		
		return this;
	},
	events: {
		'click #scene_mode': 'changeToSceneMode',
		'click #item_mode': 'changeToItemMode',
		'click #back': 'backToGameList',
		'click #save': 'save',
		'click #test_play': 'testplayHasClicked'
	},
	changeToSceneMode: function() {
		rootView.changeMode('scene_editor');
	},
	changeToItemMode: function() {
		rootView.changeMode('item_editor');
	},
	backToGameList: function() {
	},
	testplayHasClicked: function() {
		var width = 320;
		var height = 417;
		var left = window.screen.width / 2 - width / 2;
		var top = window.screen.height / 2 - height / 2;
		var options = 'width=' + width + ',height=' + height + ',left=' + left + ',top=' + top + ',status=no,location=no,toolbar=no,menubar=no';
		window.open('/runtime?game_key=' + GAME_ID, null, options);
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