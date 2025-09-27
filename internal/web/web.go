package web

import (
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	// Define headers para HTML
	headers(w)

	// Conteúdo HTML
	html := `
<!DOCTYPE html>
<html lang="pt-BR">
	<head>
		<script src="https://cdn.jsdelivr.net/npm/xterm@5.3.0/lib/xterm.min.js"></script>
		<link href="https://cdn.jsdelivr.net/npm/xterm@5.3.0/css/xterm.min.css" rel="stylesheet">
	</head>
	<body>
		<div id="terminal"></div>
		<script>
		var term = new Terminal();
		term.open(document.getElementById('terminal'));
		term.write('Hello from \x1B[1;3;31mxterm.js\x1B[0m $ ')
		</script>
	</body>
</html>
`
	w.Write([]byte(html))
}

func headers(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
