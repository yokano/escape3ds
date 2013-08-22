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
	events: {
		'change #change_scene_img': 'sceneImageHasChanged',
		'change #scene_info .scene_name': 'sceneNameHasChanged',
		'click #delete_scene': 'deleteButtonHasClicked',
		'click #scene_info .is_first_scene': 'firstSceneCheckboxHasClicked',
		'drop .dropbox': 'eventImageHasDropped',
		'dragenter .dropbox': 'eventImageHasEntered',
		'click .scene_img': 'sceneImageHasClicked',
		'click #edit_enter_event': 'editEnterEventHasClicked',
		'click #edit_leave_event': 'editLeaveEventHasClicked',
		'click #remove_scene_background': 'removeSceneBackgroundHasClicked'
	},
	
	// 描画処理
	render: function() {
		if(this.model == null) {
			this.$el.hide();
			return this;
		}
		this.$el.show();
		this.$el.html(this.template(this.model.toJSON()));
		
		// ゲーム開始時のシーンならチェック
		if(game.get('firstScene') == this.model.id) {
			this.$el.find('#scene_info .is_first_scene').attr('checked', true);
		}

		// イベントの表示
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
				if(rootView.jcropAPI.mode == 'add') {
					self.eventAreaHasSelected(event, rootView.jcropAPI, self);
				} else if(rootView.jcropAPI.mode == 'update') {
					// update
					var currentEvent = self.model.get('eventList').getSelected();
					currentEvent.set({
						position: [event.x, event.y],
						size: [event.w, event.h]
					});
					rootView.jcropAPI.mode = 'add';
				}
				rootView.jcropAPI.release();
			},
			onRelease: function(event) {
				self.model.get('eventList').trigger('cancel');
				rootView.jcropAPI.mode = 'add';
			}
		}, function() {
			rootView.jcropAPI = this;
			rootView.jcropAPI.mode = 'add'; // イベント新規追加モード
		});
		
		return this;
	},
	
	/**
	 * シーン画像（画面下）がクリックされた
	 */
	sceneImageHasClicked: function() {
		this.$el.find('#change_scene_img').click();
	},
	 
	/**
	 * シーンの画像ファイルが変更された
	 * @method
	 */
	sceneImageHasChanged: function() {
		var url = geturl();
	
		// 画像ファイルのアップロード
		var form = this.$el.find('#change_scene_img_form');
		var formData = new FormData(form.get(0));
		var self = this;
		$.ajax(url, {
			method: 'POST',
			contentType: false,
			processData: false,
			data: formData,
			dataType: 'json',
			error: function() {
				console.log('error');
			},
			success: function(data) {
				self.model.set('background', data.blobkey)
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
		game.get('sceneList').remove([scene]);
	},
		
	/**
	 * 「ゲーム開始時のシーンにする」がクリックされた
	 * @method
	 */
	firstSceneCheckboxHasClicked: function() {
		var checked = this.$el.find('#scene_info .is_first_scene:checked').length;
		if(checked == 1) {
			game.set({
				firstScene: this.model.id,
				thumbnail: this.model.get('background')
			});
		} else {
			game.set({
				firstScene: null,
				thumbnail: ''
			});
		}
	},
	
	/**
	 * シーンの情報が変更された
	 * @method
	 */
	sceneHasChanged: function() {
		var url;
		if(this.model.get('background') == '') {
			url = '/client/editor/img/black.png';
		} else {
			url = '/download?blobkey=' + this.model.get('background');
		}
		this.$el.find('#scene').css('background-image', 'url("' + url + '")');
		this.$el.find('#scene_info .scene_img').attr('src', url);
	},
	
	/**
	 * 範囲選択されたらイベントを追加する
	 * @method
	 * @param {Object} event jcropで範囲選択された時の座標を含むオブジェクト
	 * @param {Object} jcropAPI jcropの操作を行うためのオブジェクト
	 * @param {string} image 画像のURL。ドラッグアンドドロップで画像を追加した場合にセット。
	 */
	eventAreaHasSelected: function(event, jcropAPI, self, url) {
		// this は JcropObject を指す, self が view を指す
		var name = window.prompt('イベント名を入力してください');
		jcropAPI.release();
		
		if(name == '') {
			return this;
		}
		
		var e = new Event({
			name: name,
			image: (url == undefined) ? '' : url,
			code: '',
			position: [event.x, event.y],
			size: [event.w, event.h],
			sceneId: self.model.id
		});
		self.model.get('eventList').urlRoot = '/sync/event/' + self.model.id;
		self.model.get('eventList').add(e);
//		self.model.get('eventList').trigger('eventAreaHasSelected', e);

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
		var files = event.dataTransfer.files;
		var formData = new FormData();
		formData.append('file', files[0]);
		var url = geturl();

		$.ajax(url, {
			method: 'POST',
			contentType: false,
			processData: false,
			data: formData,
			dataType: 'json',
			error: function() {
				console.log('error');
			},
			success: function(data) {
				// 画像サイズを取得
				var image = new Image();
				image.src = '/download?blobkey=' + data.blobkey;
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
					self.eventAreaHasSelected(e, jcropAPIDummy, self, data.blobkey);
				};
			}
		});
		
		return false;
	},
	
	eventImageHasEntered: function(event) {
	},
	
	/**
	 * シーン開始時のイベント編集ボタンがクリックされた
	 */
	editEnterEventHasClicked: function() {
		var eventEditorWindow = window.open('/enter_event_editor?game_key=' + GAME_ID + '&scene_key=' + this.model.id, 'イベントエディタ', 'width=640, height=800px, menubar=no, location=no, status=no');
		if(eventEditorWindow == null) {
			alert('イベントエディタの起動に失敗しました。ポップアップのブロックを解除してください。');
		}
	},
	
	/**
	 * シーン終了時のイベント編集ボタンがクリックされた
	 */
	editLeaveEventHasClicked: function() {
		var eventEditorWindow = window.open('/leave_event_editor?game_key=' + GAME_ID + '&scene_key=' + this.model.id, 'イベントエディタ', 'width=640, height=800px, menubar=no, location=no, status=no');
		if(eventEditorWindow == null) {
			alert('イベントエディタの起動に失敗しました。ポップアップのブロックを解除してください。');
		}
	},
	
	/**
	 * シーンの背景画像削除ボタンがクリックされた
	 */
	removeSceneBackgroundHasClicked: function() {
		this.model.set('background', '');
		$('#change_scene_img').val('');
		return false;
	}
});