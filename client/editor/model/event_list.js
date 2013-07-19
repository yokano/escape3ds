/**
 * イベントリストモデル
 * @class
 * @extends Collection
 * @member {Object} selected 
 */
var EventList = Backbone.Collection.extend({
	model: Event,
	initialize: function() {
		this.urlRoot = '/sync/event/' + this.get('sceneId');
		this.on('eventAreaHasSelected', this.eventAreaHasSelected);
		this.on('removeButtonHasClicked', this.removeButtonHasClicked);
		this.on('add', this.eventHasAdded);
		this.on('remove', this.eventHasRemoved);
	},
	
	/**
	 * イベントが選択された
	 * @parm {Event} event 選択されたイベントモデル
	 */
	eventAreaHasSelected: function(event) {
		this.each(function(event) {
			event.set('selected', false);
		});
		event.set('selected', true);
	},
	
	/**
	 * 現在選択されているイベントを返す
	 * @returns {Event} 選択状態のイベント
	 */
	getSelected: function() {
		var selectedEvent = this.find(function(event) {
			return event.get('selected');
		});
		return selectedEvent;
	},
	
	/**
	 * イベントが追加された
	 * @param {Event} event 追加されたイベント
	 */
	eventHasAdded: function(event) {
		Backbone.sync('create', event, {
			success: function(data) {
				event.set('id', data.id);
			},
			error: function() {
				console.log('error');
			}
		});
	},
	
	/**
	 * イベントが削除された
	 * @param {Event} event 削除されたイベント
	 */
	eventHasRemoved: function(event) {
		Backbone.sync('delete', event, {
			success: function() {
				console.log('deleting event has successed');
			},
			error: function() {
				console.log('error');
			}
		});
	}
});
