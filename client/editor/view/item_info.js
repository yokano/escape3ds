/**
 * アイテム管理画面の中央下に表示されるアイテムの設定画面
 * @class
 */
var ItemInfoView = Backbone.View.extend({
	tagName: 'div',
	id: 'item_info',
	template: _.template($('#item_info_template').html()),
	initialize: function() {
		this.model = null;
		this.listenTo(game.get('itemList'), 'change:selected', function(item) {
			if(item.get('selected')) {
				this.model = item;
				this.render();
			}
		});
	},
	events: {
		'click #delete_item': 'deleteItemButtonHasClicked',
		'click #has_first': 'hasFirstCheckboxHasClicked',
		'change #item_name': 'itemNameHasChanged',
		'change #item_img_form input': 'itemImgHasChanged',
		'click .item_img': function() { $('#item_img_form input').click(); }
	},
	render: function() {
		if(this.model == null) {
			this.$el.hide();
			return this;
		}
		
		this.$el.html(this.template(this.model.toJSON()));
		
		var url = (this.model.get('img') == '') ? '/client/editor/img/blank_item.png' : '/download?blobkey=' + this.model.get('img');
		this.$el.find('.item_img').css('background-image', 'url("' + url + '")');
		this.$el.show();
		
		return this;
	},
	
	/**
	 * アイテムの削除ボタンが押された
	 */
	deleteItemButtonHasClicked: function() {
		var message = 'アイテム「' + this.model.get('name') + '」を削除しますか';
		if(window.confirm(message)) {
			game.get('itemList').remove(this.model);
		}
	},
	
	/**
	 * 「最初から持っている」がクリックされた
	 */
	hasFirstCheckboxHasClicked: function(event) {
		// 最初から持てるアイテムは 10 個まで
		if(!this.model.get('hasFirst')) {
			var hasFirstItems = game.get('itemList').filter(function(item) {
				return item.get('hasFirst');
			});
			if(hasFirstItems.length >= 10) {
				alert('最初から持てるアイテムは１０個までです。既に１０個のアイテムが設定されています。');
				return false;
			}
		}

		this.model.set('hasFirst', !this.model.get('hasFirst'));
		if(this.model.get('hasFirst')) {
			this.$el.find('has_first').attr('checked', 'true');
		} else {
			this.$el.find('has_first').attr('checked', 'false');
		}
	},
	
	/**
	 * アイテム名が変更された
	 */
	itemNameHasChanged: function(event) {
		this.model.set('name', $(event.target).val());
	},
	
	/**
	 * アイテムの画像が変更された
	 */
	itemImgHasChanged: function(event) {
		var form = this.$el.find('#item_img_form').get(0);
		var formData = new FormData(form);
		var url = geturl();
		
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
				self.model.set('img', data.blobkey)
				self.$el.find('.item_img').css('background-image', 'url("/download?blobkey=' + data.blobkey + '")');
			}
		});
	}
});