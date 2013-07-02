/**
 * イベントリストモデル
 * @class
 * @extends Collection
 * @member {Object} selected 
 */
var EventList = Backbone.Collection.extend({
	model: Event,
	initialize: function() {
		this.on('eventAreaHasSelected', this.eventAreaHasSelected);
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
		console.log(selectedEvent);
		return selectedEvent;
	}
});
