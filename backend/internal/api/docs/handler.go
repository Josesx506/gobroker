package docs

import "net/http"

func HandleDocs(w http.ResponseWriter, r *http.Request) {
	data, _ := FS.ReadFile("docs.html")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(data)
}

func HandleSpec(w http.ResponseWriter, r *http.Request) {
	data, _ := FS.ReadFile("openapi.yaml")
	w.Header().Set("Content-Type", "application/yaml")
	w.Write(data)
}
