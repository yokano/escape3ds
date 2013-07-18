/**
 * ルートビュー
 */
var RootView = Backbone.View.extend({
	id: 'root',
	tagName: 'div',
	render: function() {
		return this;
	}
});

var ItemListView = Backbone.View.extend({
	id: 'item_list',
	tagName: 'ul',
	render: function() {
		return this;
	}
});

var ItemView = Backbone.View.extend({
	className: 'item',
	tagName: 'li',
	render: function() {
		return this;
	}
});

var MessageView = Backbone.View.extend({
	id: 'message',
	tagName: 'div',
	render: function() {
		return this;
	}
});

var SceneView = Backbone.View.extend({
	id: 'scene',
	tagName: 'div',
	render: function() {
		return this;
	}
});