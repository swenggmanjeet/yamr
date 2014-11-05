package main

const index_html = `
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>yamr group chat</title>
	<link type="text/css" href="/app.css" rel="stylesheet" media="screen">
</head>
<body id="body">
	<div id="online"><div class="inner">
	</div></div>

	<div id="room" data-name="{{.Room.GetName}}" data-v="1"></div>
	
	<div id="chat"><div>
		<div id="help">?</div>
		<form id="chat-form" autocomplete="off">
			<input type="hidden" id="idle"/>
				{{if .User.IsLoggedIn}}
					<input type="file" name="headshot" class="headshot-file"/>
					{{ .HeadshotImg }}
				{{else}}
					{{ .HeadshotImg }}
				{{end}}
			<span id="user">{{ .User.Name }}:</span><input type="text" id="message"/>
			<input type="submit" class="hidden"/>
		</form>
	</div></div>

	<div id="menu">
		{{if .User.IsGuest}}
			Welcome to yamr!<hr/>
			To visit a different chat room, go to <a href="http://yamr.net/room-name">yamr.net/room-name</a><br/><br/>
			You can also <a class="signup">signup</a> or <a class="login">login</a>
		{{else}}
			{{ .User.Name }}<hr/>
			You can change your headshot by clicking on<br/>the icon in the bottom left corner<br/><br/>
			<a class="logout">Logout</a>
		{{end}}
	</div>
	
	<div class="poof"></div>
	
	<img id="spinner" src="/images/spinner.gif"/>
	
	<script src="/app.js"></script>
</body>
</html>
`
