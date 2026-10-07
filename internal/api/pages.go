package api

import (
	"log"
	"net/http"
)

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := s.templates.ExecuteTemplate(w, "upload.html", nil); err != nil {
		log.Printf("rendering upload page: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}

func (s *Server) handleUploadPage(w http.ResponseWriter, r *http.Request) {

	result, ok := s.parseUpload(w, r)
	if !ok {
		return
	}

	var encounterViews []EncounterView

	for _, encounter := range result.Encounters {
		view := toEncounterView(encounter)
		encounterViews = append(encounterViews, view)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := s.templates.ExecuteTemplate(w, "results.html", encounterViews)
	if err != nil {
		log.Printf("rendering upload page: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}
