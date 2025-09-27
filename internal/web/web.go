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
	<meta charset="UTF-8">
	<title>Minha Página</title>
</head>
<body>
	<h1>Olá, mundo!</h1>
	<p>Esta é uma página HTML servida pelo Go.</p>
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
