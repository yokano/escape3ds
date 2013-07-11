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
		this.listenTo(game.get('sceneList'), 'select', this.sceneHasSelected);
		this.listenTo(game.get('sceneList'), 'remove', this.sceneHasRemoved);
		this.children = [];
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

		var that = this;
		this.model.get('eventList').each(function(event) {
			var eventAreaView = new EventAreaView({model: event});
			that.$el.find('#scene').append(eventAreaView.render().el);
		});
		
		// Jcrop の設定
		var self = this;
		this.$el.find('#scene').Jcrop({
			bgColor: 'white',
			onSelect: function(event) {
				var jcropAPI = this;
				self.eventAreaHasSelected(event, jcropAPI, self);
			}
		});
		
		return this;
	},
	events: {
		'change #change_scene_img': 'sceneImageHasChanged',
		'change #scene_info .scene_name': 'sceneNameHasChanged',
		'click #delete_scene': 'deleteButtonHasClicked',
		'click #copy_scene': 'copyButtonHasClicked',
		'click #scene_info .is_first_scene': 'firstSceneCheckboxHasClicked',
		'drop .dropbox': 'eventImageHasDropped',
		'dragenter .dropbox': 'eventImageHasEntered'
	},
	
	/**
	 * シーンの画像ファイルが変更された
	 * @method
	 */
	sceneImageHasChanged: function() {
		var form = $('<form></form>').append($('#change_scene_img'));
		var formData = new FormData(form.get(0));
		
		// 画像ファイルのアップロード
		$.ajax('/upload', {
			method: 'POST',
			contentType: false,
			processData: false,
			data: formData,
			dataType: 'json',
			error: function() {
				console.log('error');
			},
			success: function(data) {
				console.log('blobkey', data.blobkey);
			}
		});
	},
	
	/**
	 * シーンが選択された
	 * @method
	 */
	sceneHasSelected: function(cid) {
		if(this.model != null) {
			this.stopListening(this.model.get('eventList'));
			this.stopListening(this.model);
		}
		this.model = game.get('sceneList').get(cid);
		this.listenTo(this.model, 'change', this.sceneHasChanged);				
		this.listenTo(this.model.get('eventList'), 'add', this.eventHasAdded);
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
		var scene = game.get('sceneList').get(game.get('sceneList').selected);
		if(scene == null) {
			return;
		}
		if(!window.confirm(scene.get('name') + 'を削除しますか？')) {
			return;
		}
		game.get('sceneList').remove(scene);
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
		game.get('sceneList').add(clone);
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
	 * @param {Object} event jcropで範囲選択された時の座標を含むオブジェクト
	 * @param {Object} jcropAPI jcropの操作を行うためのオブジェクト
	 */
	eventAreaHasSelected: function(event, jcropAPI, self) {
		// this は JcropObject を指す
		var name = window.prompt('イベント名を入力してください');
		jcropAPI.release();
		
		if(name == '') {
			return this;
		}
		
		var e = new Event({
			name: name,
			image: '',
			code: '',
			position: [event.x, event.y],
			size: [event.w, event.h],
			sceneKey: ''
		});
		self.model.get('eventList').add(e);
		e.trigger('eventAreaHasSelected', e);
	},
	
	/**
	 * 新しいイベントが追加されたら表示する
	 * @method
	 * @param {Event} event 追加されたイベント
	 */
	eventHasAdded: function(event) {
		var eventAreaView = new EventAreaView({model: event});
		this.$el.find('#scene').append(eventAreaView.render().el);
	},
	
	/**
	 * シーン上にイベントの画像ファイルがドロップされた
	 * @method
	 */
	eventImageHasDropped: function(event) {
		event = event.originalEvent;
		var self = this;
		
		// 画像を取得
		var file = event.dataTransfer.files[0];
		getFileURL(file, this, function(url) {
		
			// 画像サイズを取得
			var image = new Image();
			image.src = url;
			image.onload = function() {

				// jcropの範囲選択イベントになりすましてイベントを追加
				var jcropAPIDummy = {
					release: function() {}
				};
				var e = {};
				e.w = this.width;
				e.h = this.height;
				e.x = event.offsetX - e.w / 2;
				e.y = event.offsetY - e.h / 2;
				self.eventAreaHasSelected(e, jcropAPIDummy, self);
				
				// イベントの画像を設定
				self.model.get('eventList').getSelected().set('image', url);
			};
		});
		
		return false;
	},
	
	eventImageHasEntered: function(event) {
	}
});