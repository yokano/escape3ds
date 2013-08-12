/**
 * ゲーム一覧ページのスクリプト
 * @file
 */
$(function() {

	// ゲーム新規作成ボタン
	$('#add_game').click(function() {
		var div = $('#add_game_div');
		var name = div.find('.name').val();
		var description = div.find('.description').val();
		if(name == '') {
			name = '新しいゲーム'
		}
		if(description == "") {
			description = '新しいゲームです';
		}
		$.ajax('/add_game', {
			method: 'POST',
			data: {
				user_key: userKey,
				game_name: name,
				game_description: description
			},
			dataType: 'json',
			success: function(data) {
				var li = $('<li class="game" key="' + data.key + '"></li>');
				li.append('<div class="title">' + data.name + '</div>');
				li.append('<div class="description">' + data.description + '</div>');
				li.append('<div class="thumbnail"><img width="200" src="/client/gamelist/img/black.png"></div>');
				li.append('<a href="/editor?game_key=' + data.key + '"><button>作る</button></a>');
				li.append('<button class="copy">コピー</button>');
				li.append('<button class="delete">消す</button>');
				li.hide();
				$('#gamelist').append(li);
				li.ready(function() {
					li.fadeIn();
				});
			},
			error: function(xhr, err) {
				console.log(err);
			}
		});
	});
	
	// ゲーム削除ボタン
	$(document).on('click', '.game .delete', function() {
		if(!window.confirm('ゲームを削除しますか？')) {
			return false;
		}
		var game = $(this).parent('.game');
		$.ajax('/delete_game', {
			method: 'POST',
			dataType: 'json',
			data: {
				game_key: game.attr('key')
			},
			error: function() {
				console.log('delete game error');
			},
			success: function() {
				game.fadeOut('fast');
			}
		});
	});
});