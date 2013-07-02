/**
 * シーンの設定ビュー
 * 画面中央
 * @class
 * @extends Backbone.View
 */
var SceneView = Backbone.View.extend({
	tagName: 'div',
	id: 'scene_view',
	template: _.template($('#scene_view_template').html()),
	model: null,
	initialize: function() {
		this.listenTo(sceneList, 'select', this.sceneHasSelected);
		this.listenTo(sceneList, 'remove', this.sceneHasRemoved);
		this.listenTo(eventList, 'add', this.eventHasAdded);
	},
	render: function() {
		if(this.model == null) {
			this.$el.hide();
			return this;
		}
		this.$el.show();
		this.$el.html(this.template(this.model.toJSON()));
		if(game.get('firstScene') == this.model) {
			this.$el.find('#scene_info .is_first_scene').attr('checked', true);
		}
		
		// Jcrop の設定
		var jcropAPI;
		this.$el.find('#scene').Jcrop({
			bgColor: 'white',
			onSelect: this.eventAreaHasSelected
		});
		
		return this;
	},
	events: {
		'change #change_scene_img': 'sceneImageHasChanged',
		'change #scene_info .scene_name': 'sceneNameHasChanged',
		'click #delete_scene': 'deleteButtonHasClicked',
		'click #copy_scene': 'copyButtonHasClicked',
		'click #scene_info .is_first_scene': 'firstSceneCheckboxHasClicked'
	},
	
	/**
	 * シーンの画像ファイルが変更された
	 * @method
	 */
	sceneImageHasChanged: function() {
		var that = this;
		var file = $('#change_scene_img').get(0).files[0];
		var fileReader = new FileReader();
		fileReader.onload = function(data) {
			that.model.set('background', data.target.result);
		};
		fileReader.readAsDataURL(file);

//		保存処理へ移動させる
//		var form = $('#change_scene_img_form').get()[0];
//		var formData = new FormData(form);
//		$.ajax('/upload', {
//			method: 'POST',
//			contentType: false,
//			processData: false,
//			data: formData,
//			dataType: 'json',
//			error: function() {
//				console.log('error');
//			},er
//			success: function(data) {
//				console.log('blobkey', data.blobkey);
//			}
//		});
	},
	
	/**
	 * シーンが選択された
	 * @method
	 */
	sceneHasSelected: function(cid) {
		this.model = sceneList.get(cid);
		this.listenTo(this.model, 'change', this.sceneHasChanged);
		this.render();
	},
	
	/**
	 * シーン名が変更された
	 * @method
	 */
	sceneNameHasChanged: function() {
		var name = this.$el.find('#scene_info .scene_name').val();
		this.model.set('name', name);
	},
	
	/**
	 * シーンが削除された
	 * @method
	 */
	sceneHasRemoved: function() {
		this.stopListening(this.model);
		this.model = null;
		this.render();
	},
	
	/**
	 * 「シーンを削除」ボタンがクリックされた
	 * @method
	 */
	deleteButtonHasClicked: function() {
		var scene = sceneList.get(sceneList.selected);
		if(scene == null) {
			return;
		}
		if(!window.confirm(scene.get('name') + 'を削除しますか？')) {
			return;
		}
		sceneList.remove(scene);
	},
	
	/**
	 * 「シーンをコピー」ボタンがクリックされた
	 * @method
	 */
	copyButtonHasClicked: function() {
		var name = window.prompt('コピー先のシーン名を入力してください');
		if(name == '') {
			return;
		}
		var clone = this.model.clone();
		clone.set('name', name);
		sceneList.add(clone);
	},
	
	/**
	 * 「ゲーム開始時のシーンにする」がクリックされた
	 * @method
	 */
	firstSceneCheckboxHasClicked: function() {
		var checked = this.$el.find('#scene_info .is_first_scene:checked').length;
		if(checked == 1) {
			game.set('firstScene', this.model);
		} else {
			game.set('firstScene', null);
		}
	},
	
	/**
	 * シーンの情報が変更された
	 * @method
	 */
	sceneHasChanged: function() {
		this.$el.find('#scene').css('background-image', 'url(' + this.model.get('background') + ')');
		this.$el.find('#scene_info .scene_img').attr('src', this.model.get('background'));
	},
	
	/**
	 * 範囲選択されたらイベントを追加する
	 * @method
	 * @param {Object} jcropAPI
	 */
	eventAreaHasSelected: function(data) {
		// this は JcropObject を指す
		var name = window.prompt('イベント名を入力してください');
		this.release();
		
		if(name == '') {
			return;
		}
		
		eventList.add({
			name: name,
			image: '',
			code: '',
			position: [data.x, data.y],
			size: [data.w, data.h],
			sceneKey: ''
		});
	},
	
	/**
	 * 新しいイベントが追加されたら表示する
	 * @method
	 * @param {Event} 追加されたイベント
	 */
	eventHasAdded: function(data) {
		var eventAreaView = new EventAreaView();
	}
});