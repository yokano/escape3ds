/**
 * シーン一覧の１つ１つのli要素
 * @class
 * @extends Backbone.View
 * @member {SceneListView} parent 親要素への参照
 */
var SceneListItemView = Backbone.View.extend({
	tagName: 'li',
	template: _.template($('#scene_li_view_template').html()),
	parent: null,
	initialize: function(config) {
		this.parent = config.parent;
		this.listenTo(this.model, 'change', this.sceneHasChanged);
		this.listenTo(this.model, 'remove', this.sceneHasRemoved);
	},
	render: function() {
		this.$el.html(this.template(this.model.toJSON()));
		return this;
	},
	events: {
		'click': 'selectScene'
	},
	selectScene: function(event) {
		if(this.$el.hasClass('select')) {
			return;
		}
		this.parent.collection.trigger('select', this.model.cid)
	},
	sceneHasChanged: function() {
		this.$el.find('.scene_name').html(this.model.get('name'));
		this.$el.find('.scene_img').attr('src', this.model.get('background'));
	},
	sceneHasRemoved: function() {
		this.remove();
	}
});
