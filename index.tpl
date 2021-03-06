<!DOCTYPE html>
<html>
	<head>
		<title>{{.Title}}</title>
		<style>
			body {
				font-family: monospace;
			}
			#screen {
				border: 2px outset black;
				border-radius: 5px;
			}
			#logo {
				height: 50px;
				width: auto;
				float: right;
				margin-right: 10px;
			}
		</style>
	</head>
	<body>
		<h1>{{.Title}}<img id="logo" src="/static/nullptr.png" /></h1>
		<img id="screen" src="/img" alt="Simulation Screen"/>
		<script>
			/*
			ws = new WebSocket("{{.Ws}}");
			ws.onopen = function(e) {
				console.log("OPEN");
			}
			ws.onclose = function(e) {
				console.log("CLOSE");
				ws = null;
			}
			ws.onmessage = function(e) {
				console.log("RESPONSE");

				var canvas = document.createElement('canvas');
				canvas.height = {{.Height}};
				canvas.width = {{.Width}};
				var ctx = canvas.getContext('2d');

				document.getElementById("screen").remove()
				document.getElementsByTagName("body")[0].insertAdjacentHTML("beforeend", "<img id=\"screen\" src=\"data:image/png;base64, " + e.data + "\" alt=\"Simulation Screen\"/>")
			}
			ws.onerror = function(e) {
				console.log("ERROR");
			}*/
			setInterval(refresh,1000)
			function refresh() {
				document.getElementById("screen").src =
					"/img?t=" + new Date().getTime();
			}
		</script>
	</body>
</html>
