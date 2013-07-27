var BlockList = Backbone.Collection.extend({
});

var MethodBlock = Backbone.Model.extend({
	defaults: {
		type: '',
		attr: ''
	}
});

var IfBlock = Backbone.Model.extend({
	defaults: {
		type: '',
		attr: '',
		true: null,
		false: null
	}
});