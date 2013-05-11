/**
 * エディタのスクリプト
 * @file
 */
$(function() {
	if(!window.File) {
		alert('お使いのブラウザは File API に対応していません\n最新のブラウザをご使用ください');
		return false;
	}
	
	$('#scene').on('drop', function(event) {
		var file = event.originalEvent.dataTransfer.files[0];
		if(!file.type.match(/image*/)) {
			alert('画像ファイルをアップロードしてください');
			return false;
		}
		
		var reader = new FileReader();
		reader.onerror = function() {
			console.log('file read error');
		};
		reader.onload = function(event) {
			$('#scene').css('background-image', 'url(' + event.target.result + ')');
		}
		reader.readAsDataURL(file);
		$(this).css('background-color', 'black');
		
		var jcropApi;
		$(this).Jcrop({
			onSelect: function(data) {
				$('<div class="event_area"></div>')
					.css('left', data.x)
					.css('top', data.y)
					.css('width', data.w)
					.css('height', data.h)
					.appendTo($('#scene'));
				jcropApi.release();
			}
		}, function() {
			jcropApi = this;
		});
		
		return false;

	}).on('dragenter', function() {
		$(this).empty().css('background-color', 'gray');
	}).on('dragleave', function() {
		$(this).html('<br/><br/><br/><br/>背景画像ファイルをここへドラッグ<br/>サイズは横320x212pxです').css('background-color', 'black');
	});
});