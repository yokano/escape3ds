/**
 * ログインページのスクリプト
 * @file
 */
$(function() {
	// ログイン処理
	var loginForm = $('.normal_login');
	loginForm.find('.submit').click(function() {
		var mail = loginForm.find('.mail').val();
		var pass = loginForm.find('.pass').val();
		$.ajax('/login', {
			method: 'POST',
			data: {
				mail: mail,
				pass: pass
			},
			dataType: 'json',
			success: function(data) {
				if(data.result == false) {
					alert(data.message);
				} else {
					location.href = data.to;
				}
			},
			error: function() {
				console.log('error');
			}
		});
	});
	
	// パスワード忘れたボタン
	$('#forget').on('click', function() {
		var confirm = window.confirm('パスワードを再発行しますか？');
		if(!confirm) {
			return;
		}
		
		var mail = window.prompt('登録したメールアドレスを入力してください');
		if(mail == '') {
			return;
		}
		
		$.ajax('/reset_password', {
			method: 'POST',
			data: {
				mail: mail
			},
			success: function() {
				alert('新しいパスワードをメールアドレスへ送信しました。ご確認ください。');
			},
			error: function() {
				alert('パスワードの再発行に失敗しました。正しいメールアドレスを入力したかご確認ください。');
			}
		});
	});
	
	// 新規登録
	var registration = $('.registration');
	registration.find('.submit').on('click', function() {
		var name = registration.find('.name').val();
		var mail = registration.find('.mail').val();
		var pass = registration.find('.password').val();
		$.ajax('/interim_registration', {
			method: 'POST',
			data: {
				name: name,
				mail: mail,
				pass: pass
			},
			dataType: 'json',
			success: function(data) {
				if(data.result == false) {
					alert(data.msg);
					return;
				}
				location.href = "/registration_successed";
			},
			error: function(data) {
				console.log('registration error');
			}
		});
	});
});