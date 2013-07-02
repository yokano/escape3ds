/**
 * イベントリストモデル
 * @class
 * @extends Collection
 * @member {Object} selected 
 */
var EventList = Backbone.Collection.extend({
	model: Event,
	selected: null
});