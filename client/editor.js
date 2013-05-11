/**
 * エディタのスクリプト
 * @file
 */
$(function() {
	// 設定ファイル作成
	var config = {
		title: 'サンプル',
		start: 'default',
		author: 'y.okano',
		scenes: [
			{
				id: 'default',
				events: []
			}
		]
	};
	var currentScene = 'default'
	
	// シーンリストを選択可能にする
	$(document).on('click', '#scene_list>li', function() {
		$('#scene_list>li').removeClass('selected');
		$(this).addClass('selected');
	});
	
	// イベント領域がクリックされたらイベントを開く
	$(document).on('click', '.event_area:not(.selected)', function() {
		$('.event_area').removeClass('selected');
		$(this).addClass('selected');
		console.log('イベントが選択されました');
	});
	
	// シーン追加ボタン
	$('#add_scene').on('click', function() {
		$('#scene_list').append($('<li>新しいシーン</li>'));
	});
	
	// File API 対応チェック
	if(!window.File) {
		alert('お使いのブラウザは File API に対応していません\n最新のブラウザをご使用ください');
		return false;
	}
	
	// 画像をシーンへドロップしたら
	$('#scene').on('drop', function(event) {
		var file = event.originalEvent.dataTransfer.files[0];
		if(!file.type.match(/image*/)) {
			alert('画像ファイルをアップロードしてください');
			return false;
		}
		
		// 画像ファイル読み込み
		var reader = new FileReader();
		reader.onerror = function() {
			console.log('file read error');
		};
		reader.onload = function(event) {
			$('#scene').css('background-image', 'url(' + event.target.result + ')');
			$('#scene_information').css('display', 'block');
		}
		reader.readAsDataURL(file);
		$(this).css('background-color', 'black');
		
		// 範囲指定イベント
		var jcropApi;
		$(this).Jcrop({
			onSelect: function(data) {
				$('<div class="event_area"></div>')
					.css('left', data.x)
					.css('top', data.y)
					.css('width', data.w)
					.css('height', data.h)
					.on('click', function() {
						console.log('ok');
					})
					.appendTo($('#scene'));
				jcropApi.release();
			}
		}, function() {
			jcropApi = this;
		});
		
		// ページ遷移をキャンセル
		return false;

	}).on('dragenter', function() {
		$(this).empty().css('background-color', 'gray');
	}).on('dragleave', function() {
		$(this).html('<br/><br/><br/><br/>背景画像ファイルをここへドラッグ<br/>サイズは横320x212pxです').css('background-color', 'black');
	});
});