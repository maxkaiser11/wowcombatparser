package api

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"github.com/maxkaiser11/wowcombatparser/internal/web"
)

type Server struct {
	templates *template.Template
}

func NewServer() (*Server, error) {
	funcs := template.FuncMap{
		"add": func(a, b int) int { return a + b },
	}
	tmpl, err := template.New("").Funcs(funcs).ParseFS(web.Templates, "templates/*.html")

	if err != nil {
		return nil, fmt.Errorf("parsing templates: %w", err)
	}
	return &Server{templates: tmpl}, nil
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.handleHealth)
	mux.HandleFunc("GET /{$}", s.handleIndex)
	mux.HandleFunc("POST /upload", s.handleUploadPage)

	static, err := fs.Sub(web.Static, "static")
	if err != nil {
		log.Fatal(err)
	}
	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.FS(static))))
	return mux
}
