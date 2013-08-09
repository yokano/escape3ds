/**
 * 汎用関数群
 */

// ファイル読み込み関数
var getFileURL = function(file, caller, callback) {
	var fileReader = new FileReader();
	fileReader.onload = function(data) {
		callback.call(caller, data.target.result);
	};
	fileReader.readAsDataURL(file);
};

// blobsotre へファイルをアップロードするための URL をサーバから発行してもらう
var geturl = function() {
	var url = '';
	
	$.ajax('/geturl', {
		method: 'GET',
		dataType: 'json',
		async: false,
		error: function() {
			console.log('アップロード先URLの取得に失敗しました');
		},
		success: function(data) {
			url = data.url;
		}
	});
	
	return url;
};
