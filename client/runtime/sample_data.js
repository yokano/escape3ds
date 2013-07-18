var game = {
	name: 'sample game',
	description: 'this is sample game.',
	first_scene: 'bedroom',
	scene_list: {
		'bedroom1': {
			name: 'bedroom1',
			background: '/client/runtime/bedroom1.png',
			enter: {
			
			},
			leave: {
			
			},
			event_list: {
				'event1': {
					name: 'event1',
					img: '',
					position: [120, 50],
					size: [30, 20],
					code: 'alert("hello");'
				}
			}
		}
	},
	item_list: {
		'hammer': {
			name: '金槌',
			has_first: true,
			img: '/client/runtime/hammer.png'
		},
		'dish': {
			name: '皿',
			has_first: false,
			img: '/client/runtime/dish.png'
		}
	}
};