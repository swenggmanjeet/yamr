var chat = {
	enabled: true,
	roomHeight: 0,
	socket: false,
	reconnect: false,
	idle: 0,
	idleTimer: false,
	typing: 0,
	paused: false,
	queue: [],

	init: function() {
		chat.disable();
		chat.resize();
		chat.connect();
		chat.setIdle();

		$('#chat-form').on('submit', function() {
			if ('/upload' !== $(this).attr('action')) {
				chat.save();
				return false;
			}
		})

		$('#message').keyup(function() {
			chat.setTyping();
		})

		document.title = $('#room').data('name');

		window.onresize = function() {
			chat.resize();
		}

		document.onmousemove = function() {
			chat.setIdle();
			document.title = $('#room').data('name');
		}

		document.onkeypress = function() {
			chat.setIdle();
			document.title = $('#room').data('name');
		}

		$(document).on('mouseover mouseout', 'div.message div div', function(e) {
			if (e.type == 'mouseout') {
				$(this).find('span.time').hide();
			} else if (e.type == 'mouseover') {
				$(this).find('span.time').show();
			}
		})

		$('.headshot-file').on('mouseover mouseout', function(e) {
			if (e.type == 'mouseout') {
				$(this).parent().find('img').css('border-color', '#262626');
			} else if (e.type == 'mouseover') {
				$(this).parent().find('img').css('border-color', '#fff');
			}
		})

		// $('#hiddenIframe').load(chat.hiddenFrameLoad);
		// $('input.headshot-file').on('change', function() {
		//  $('form#headshot-form').submit();
		//  chat.showInfo('<img src="/images/spinner.gif" width="14" height="14"/>');
		// })
		$('input.headshot-file').on('change', chat.uploadHeadshot);

		$(document).click(chat.handleClick);
	},

	uploadHeadshot: function() {
		chat.showInfo('<img src="/images/spinner.gif" width="14" height="14"/>');

		var form = $('#chat-form');
		var iframe = document.createElement('iframe');
		iframe.setAttribute('id', 'hiddenIframe');
		iframe.setAttribute('name', 'hiddenIframe');
		iframe.setAttribute('class', 'hidden');

		form.before(iframe);

		$(iframe).load(chat.hiddenFrameLoad);

		window.frames['hiddenIframe'].name = 'hiddenIframe';

		form.attr('target', 'hiddenIframe');
		form.attr('action', '/upload');
		form.attr('method', 'POST');
		form.attr('enctype', 'multipart/form-data');
		form.submit();

		form.attr('action', '');
	},

	hiddenFrameLoad: function() {
		if (window.frames['hiddenIframe'].document == null) {
			chat.showInfo('error adding headshot');
			return;
		}

		var iframe = frames['hiddenIframe'].document.body.innerText;

		if (iframe.indexOf('/headshots/') !== -1) {
			chat.showInfo('?');
			user = chat.getCurrentUser();
			$('#u' + user + ' img').attr('src', iframe);
			$('#chat').find('img').attr('src', iframe);
		} else {
			if (!iframe) {
				iframe = 'error adding headshot'
			}
			chat.showInfo(iframe);
		}
	},

	scrollToBottom: function(id) {
		var id = $.trim(id);
		document.getElementById(id).scrollTop = document.getElementById(id).scrollHeight;
	},

	enable: function() {
		chat.enabled = true;
		$('#body').removeClass('faded');
		$('#spinner').hide();
		$('#message').attr('disabled', false);
		$('#message').focus();
		chat.showInfo('?');
	},

	disable: function(msg) {
		var w = $('#body').width();
		var h = $('#body').height();

		if ('undefined' === typeof(msg)) {
			$('#spinner').css({
				position: 'absolute',
				left: (w/2) - 25,
				top: (h/2) - 25
			}).show();
		} else {
			var div = $('<div class="hidden">'+msg+'</div>');
			$('#body').append(div);
			var dw = div.width();
			var dh = div.height();
			
			$('#spinner').hide();

			div.css({
				position: 'absolute',
				left: (w/2) - (dw/2),
				top: (h/2) - (dh/2),
				'text-align': 'center'
			}).show()
		}

		$('#body').addClass('faded');
		$('#message').attr('disabled', true);
		chat.showInfo('');
	},

	resize: function() {
		var w = $('#body').width();
		var h = $('#body').height();
		var chatHeight = $('#chat').height();
		chat.roomHeight = h - chatHeight;
		$('#room').css('width', w - 210);
		$('#room').css('bottom', chatHeight);
		chat.resizeMessageInput();
		chat.scrollToBottom('room');
	},

	resizeMessageInput: function() {
		var chatWidth = $('#chat div').width();
		chatWidth -= ($('#chat div form span#user').width() + 12);
		chatWidth -= $('#help').width(); // help
		chatWidth -= 47; // headshot
		chatWidth -= 22; // misc padding
		$('#message').css('width', chatWidth);
	},

	showInfo: function(info) {
		chat.hideHelpMenu();

		$('#help').html(info);
		chat.resizeMessageInput();
	},

	toggleHelpMenu: function() {
		if ($('#menu').is(':visible')) {
			chat.hideHelpMenu();
		} else {
			chat.showHelpMenu();
		}
	},

	showHelpMenu: function() {
		$('#menu').show();
		$('#help').html('?').addClass('active');
	},

	hideHelpMenu: function() {
		$('#menu').hide();
		$('#help').removeClass('active');
	},

	playYoutube: function(obj) {
		var scrollBottom = chat.isScrolledBottom('room');
		var v = $(obj).data('v');
		$(obj).wrap('<span>').parent().html('<iframe width="281" height="200" src="http://www.youtube.com/embed/' + v + '?autoplay=1" frameborder="0"></iframe>');
		if (scrollBottom) {
			chat.scrollToBottom('room');
		}
	},

	autoLink: function(t) {
		t = ' ' + t;
		regexp = new RegExp('(http://|https://)?(www.)?youtube.com/watch\\?v=([\-_a-z0-9]+)\\S*', 'i');
		t = t.replace(regexp, '<img src="http://i1.ytimg.com/vi/$3/default.jpg" width="120" height="90" class="youtube" data-v="$3"/>');

		regexp = new RegExp('(^|[\\n ])([\\w]+?://[\\w]+[^ \\"\\n\\r\\t<]*)', 'i');
		t = t.replace(regexp, '$1<a href="$2" target="_blank" rel="nofollow">$2</a>');

		regexp = new RegExp('(^|[\\n ])(www\\.[^ \\"\\t\\n\\r<]*)', 'i');
		t = t.replace(regexp, '$1<a href="http://$2" target="_blank" rel="nofollow">$2</a>');

		regexp = new RegExp('(^|[\\n ])([a-z0-9&\\-_\\.]+?)@([\\w\\-]+\\.([\\w\\-\\.]+\\.)*[\\w]+)', 'i');
		t = t.replace(regexp, '$1<a href="mailto:$2@$3">$2@$3</a>');

		return $.trim(t);
	},

	escapeHtml: function(t) {
		var r1 = new RegExp('<', 'g');
		var r2 = new RegExp('>', 'g');
		return t.replace(r1, '&lt;').replace(r2, '&gt;');
	},

	signup: function(str) {
		var space = str.indexOf(' ');

		if (space === -1) {
			space = str.length;
		}

		chat.showInfo('<img src="/images/spinner.gif" width="14" height="14"/>');

		var args = {
			username: $.trim(str.substr(0, space)),
			password: $.trim(str.substr(space + 1))
		}

		if (!args.username) {
			chat.showInfo('Please enter a username');
			return;
		} else if (!args.password) {
			chat.showInfo('Please enter a password');
			return;
		}

		$.ajax({
			url: '/signup',
			type: 'POST',
			data: args,
			success: function(resp) {
				if ('ok' === resp) {
					location.href = location.href;
				} else {
					if (!resp) {
						resp = 'Unknown error';
					}

					chat.showInfo(resp);
				}
			}
		})
	},

	login: function(str) {
		var space = str.indexOf(' ');

		if (space === -1) {
			space = str.length;
		}

		chat.showInfo('<img src="/images/spinner.gif" width="14" height="14"/>');

		var args = {
			username: $.trim(str.substr(0, space)),
			password: $.trim(str.substr(space + 1))
		}

		if ('' === args.username || '' === args.password) {
			chat.showInfo('Invalid login');
			return;
		}

		$.ajax({
			url: '/login',
			type: 'POST',
			data: args,
			success: function(resp) {
				if ('ok' === resp) {
					location.href = location.href;
				} else {
					if (!resp) {
						resp = 'Unknown error';
					}

					chat.showInfo(resp);
				}
			}
		})
	},

	save: function() {
		var message = $.trim($('#message').val());
		var time = chat.getNow();

		$('#message').val(''); // clear ASAP to avoid double clicks

		if (!message) {
			return;
		}

		if (message === '/clear') {
			chat.clear();
			chat.showInfo('?');
			return;
		}

		if (message.substr(0, 6) === '/login') {
			chat.login(message.substr(7));
			return;
		}

		if (message.substr(0, 7) === '/signup') {
			chat.signup(message.substr(8));
			return;
		}

		if (message.substr(0, 9) === '/headshot') {
			chat.headshot(message.substr(10));
			return;
		}

		if (message === '/logout') {
			chat.logout();
			return;
		}

		chat.send({ message: message });
		chat.setIdle();
	},

	getNow: function() {
		return chat.formatTime();
	},

	formatTime: function(date) {
		if ((date instanceof Date) === false) {
			date = new Date();
		}

		var h = date.getHours();
		var i = date.getMinutes();

		if (i < 10) {
			i = '0' + i;
		}

		if (h > 12) {
			time = (h-12) + ':' + i;
		} else if (h == 0) {
			time = '12:' + i;
		} else {
			time = h + ':' + i;
		}

		return h >= 12 ? time + ' pm' : time + ' am';
	},

	clear: function() {
		$('#room').html('');
		$('#message').val('');
	},

	autoClear: function() {
		if ($('div.message').length <= 50) {
			return false;
		}

		$('div.message:first').remove();
	},

	pauseQueue: function() {
		chat.paused = true;
	},

	resumeQueue: function() {
		chat.paused = false;
		setTimeout(chat.runQueue, 1);
	},

	runQueue: function() {
		if (!chat.paused && chat.queue.length) {
			chat.queue.shift()();
			if (!chat.paused) {
				chat.resumeQueue();
			}
		}
	},

	showJoin: function(obj) {
		chat.pauseQueue();

		var user = obj.user;
		var imgId = 'u' + user.name;
		var headshot;

		if (user.headshot > 0) {
			headshot = '/headshots/' + user.name + '.jpg?' + user.headshot;
		} else {
			headshot = '/images/no_photo.gif';
		}

		if (0 === $('#' + imgId).length) {
			var html = '<div id="' + imgId + '" data-count="1">';
			html += '<img src="' + headshot + '" class="headshot" title="' + user.name + ' is online"/>';
			html += user.name;
			html += '</div>';

			$('#online').append(html);

			if (user.typing) {
				$('div#' + imgId + ' img').addClass('typing').attr('title', user.name + ' is typing');
			}

			if (user.idle) {
				$('div#' + imgId + ' img').addClass('idle').attr('title', user.name + ' is idle');
			}
		} else {
			// if user logs in with a new tab, increment the count
			// so we know not to remove the user when the tab closes
			// if the count is greater than 1
			$('#' + imgId).data('count', parseInt($('#' + imgId).data('count')) + 1);
		}

		chat.resumeQueue();
	},

	zeroPad: function(number, width) {
		width -= number.toString().length;
		if (width > 0) {
			return new Array( width + (/\./.test( number ) ? 2 : 1) ).join( '0' ) + number;
		}
		return number + ''; // always return a string
	},

	showMessage: function(obj) {
		// if (parseInt($('#room').attr('data-v')) != parseInt(obj.user.version)) {
		//  if (!chat.isTyping()) {
		//    location.href = location.href;
		//    return;
		//  }
		// }

		var lastUser = $('#room img.headshot:last');
		var scrolledBottom = chat.isScrolledBottom('room');
		var headshot;

		if (obj.user.headshot > 0) {
			headshot = '/headshots/' + obj.user.name + '.jpg?' + obj.user.headshot;
		} else {
			headshot = '/images/no_photo.gif';
		}

		var date = new Date();

		if ('created' in obj) {
			var MS_PER_MINUTE = 60000;
			// var utc = new Date(obj.created);
			// var d = new Date(utc.valueOf() -  date.getTimezoneOffset() * MS_PER_MINUTE);
			var d = new Date(obj.created)
			var time = chat.formatTime(d);
			var scrolledBottom = true;

			if (date.toDateString() == d.toDateString()) {

			} else if (date.getFullYear() == d.getFullYear()) {
				time = (d.getMonth() + 1) + '/' + d.getDate() + ' @ ' + time;
			} else {
				time = (d.getMonth() + 1) + '/' + d.getDate() + '/' + d.getFullYear() + ' @ ' + time;
			}
		} else {
			var time = chat.formatTime(date)
		}

		if (lastUser.attr('title') == obj.user.name) {
			var html = '<div>';
			html += chat.autoLink(chat.escapeHtml(obj.message));
			html += '<span class="time">'+time+'</span>';
			html += '</div>';
			lastUser.parent().children('div').append(html);
		} else {
			var html = '<div class="message clearfix">';
			html += '<img src="' + headshot + '" class="headshot" title="' + obj.user.name + '"/>';
			html += '<div>';
			html += '<strong>' + obj.user.name + ':</strong><br/>';
			html += '<div>';
			html += '<div>';
			html += chat.autoLink(chat.escapeHtml(obj.message));
			html += '<span class="time">' + time + '</span>';
			html += '</div>';
			html += '</div>';
			html += '</div>';
			$('#room').append(html);
		}

		if (chat.roomHeight && $('#room').height() > chat.roomHeight) {
			$('#room').css('height', chat.roomHeight);
			chat.roomHeight = 0;
		}

		if (scrolledBottom) {
			chat.scrollToBottom('room');
		}

		user = chat.getCurrentUser();
		if (user != obj.user.name) {
			document.title = obj.user.name + ' says...';
		}

		chat.autoClear();
	},

	getCurrentUser: function() {
		if (!chat.user) {
			chat.user = $('#chat form span#user').html().split(':')[0];
		}

		return chat.user;
	},

	isScrolledBottom: function(id) {
		var currentHeight = 0;
		var scrollHeight = document.getElementById(id).scrollHeight;
		var offsetHeight = document.getElementById(id).offsetHeight;
		var scrollTop = document.getElementById(id).scrollTop;
		var pixelHeight = document.getElementById(id).style.pixelHeight;

		if (typeof pixelHeight === 'undefined') {
			pixelHeight = 0;
		}

		if (scrollHeight > 0) {
			currentHeight = scrollHeight;
		} else if (offsetHeight > 0) {
			currentHeight = offsetHeight;
		}

		if (pixelHeight > 0) {
			offsetHeight = pixelHeight
		}

		return (currentHeight - scrollTop - offsetHeight < 50);
	},

	setIdle: function() {
		var d = new Date();
		$('#idle').val(d.getTime());

		if (chat.idleTimer) {
			clearTimeout(chat.idleTimer);
		}

		chat.isIdle();
	},

	isIdle: function() {
		var d = new Date();
		var i = $('#idle').val();
		var idle = 0;

		if (i === '') {
			i = d.getTime();
		}

		// 5 minutes
		if ((d.getTime() - i) >= 30000) {
			idle = 1;
		}

		if (idle != chat.idle) {
			chat.idle = idle;

			if (idle) {
				chat.send({ action: 'idle' });
			} else {
				chat.send({ action: '-idle' });
			}
		}

		chat.idleTimer = setTimeout(chat.isIdle, 60000);
	},

	showIdle: function(obj) {
		if ('idle' === obj.action) {
			$('#u' + obj.user.name + ' img').addClass('idle').attr('title', obj.user.name + ' is idle');
		} else {
			$('#u' + obj.user.name + ' img').removeClass('idle').attr('title', obj.user.name + ' is online');
		}
	},

	showTyping: function(obj) {
		if ('typing' == obj.action) {
			$('#u' + obj.user.name + ' img').addClass('typing').attr('title', obj.user.name + ' is typing');
		} else {
			$('#u' + obj.user.name + ' img').removeClass('typing').attr('title', obj.user.name + ' is online');
		}
	},

	setTyping: function(obj) {
		var typing = chat.isTyping();

		if (typing != chat.typing) {
			chat.typing = typing;

			if (typing) {
				chat.send({ action: 'typing' });
			} else {
				chat.send({ action: '-typing' });
			}
		}
	},

	isTyping: function() {
		var message = $('#message').val();

		if (message && message.substr(0, 1) !== '/') {
			return 1;
		} else {
			return 0;
		}
	},

	logout: function() {
		$.post('/logout', {}, function() {
			location.href = '/';
		})
	},

	handleClick: function(event) {
		if (!chat.enabled) {
			return false;
		}

		var target = $(event.target);

		if (target.is('div#help')) {
			chat.toggleHelpMenu()
		} else if (target.is('img.youtube')) {
			chat.playYoutube(target);
		} else if (target.is('a.logout')) {
			$('#message').val('/logout').focus();
			chat.hideHelpMenu();
		} else if (target.is('a.signup')) {
			$('#message').val('/signup username password').focus();
			chat.hideHelpMenu();
		} else if (target.is('a.login')) {
			$('#message').val('/login username password').focus();
			chat.hideHelpMenu();
		}
	},

	showLogout: function(user) {
		chat.pauseQueue();

		var obj = $('#u' + user.name);
		var count = parseInt(obj.data('count'));

		if (count > 1) {
			obj.data('count', count - 1);
			chat.resumeQueue();
		} else {
			obj.data('count', null);

			var offset = obj.offset();

			obj.fadeOut('fast', function() {
				$(this).remove();
				chat.resumeQueue();
			});

			$('<div class="poof">').css({
				left: offset.left + 5,
				top: offset.top + 5
			}).appendTo('body').show();

			chat.animatePoof();
		}
	},

	animatePoof: function() {
		var bgTop = 0;
		var frames = 5;
		var frameSize = 32;
		var frameRate = 80;

		for (var i = 0; i <= frames; i++) {
			$('div.poof').animate({
				backgroundPosition: '0' + (bgTop - frameSize)
			}, frameRate);
			bgTop -= frameSize;
		}

		setTimeout("$('div.poof').remove()", frames + frameRate);
	},

	handleObj: function(obj) {
		if ('message' in obj) {
			chat.showMessage(obj);
		} else if ('join' === obj.action) {
			chat.queue.push(function() { chat.showJoin(obj) });
			chat.runQueue();
		} else if ('typing' === obj.action || '-typing' === obj.action) {
			chat.showTyping(obj);
		} else if ('idle' === obj.action || '-idle' === obj.action) {
			chat.showIdle(obj);
		} else if ('logout' === obj.action) {
			chat.queue.push(function() { chat.showLogout(obj.user) });
			chat.runQueue();
		}
	},

	send: function(obj) {
		obj.room = {
			name: $('#room').data('name')
		}
		chat.socket.send(JSON.stringify(obj));
	},

	connect: function() {
		if (chat.reconnect) {
			clearInterval(chat.reconnect);
		}
		
		if (window['WebSocket']) {
			// location.pathname
			var ws_url = 'ws://'+location.host+':8000/ws?r='+$('#room').data('name');
			chat.socket = new WebSocket(ws_url);

			chat.socket.onclose = function(e) {
				chat.disable();
				chat.reconnect = setInterval(chat.connect, 10000);
			}

			chat.socket.onmessage = function(e) {
				var obj = $.parseJSON(e.data);
				chat.handleObj(obj);
			}

			chat.socket.onopen = function() {
				if (chat.reconnect) {
					clearInterval(chat.reconnect);
				}

				chat.send({ action: 'join' });

				chat.enable();
			}

			chat.socket.onerror = function(e) {
				chat.disable();
				chat.reconnect = setInterval(chat.connect, 10000);
			}
		} else {
			var msg = 'Sorry, your browser does not have WebSocket support.<br/>';
			msg += 'Please try using the latest version of one of the following browsers<br/><br/>';
			msg += '<a href="http://google.com/chrome/" class="chrome"></a>';
			msg += '<a href="http://www.apple.com/safari/" class="safari"></a>';
			msg += '<a href="http://getfirefox.com/" class="firefox"></a>';

			chat.disable(msg);
		}
	}
};