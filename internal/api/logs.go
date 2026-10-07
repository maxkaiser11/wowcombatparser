package api

import (
	"errors"
	"log"
	"net/http"

	"github.com/maxkaiser11/wowcombatparser/internal/combatlog"
)

const maxUploadSize = 1 << 30

func (s *Server) parseUpload(w http.ResponseWriter, r *http.Request) (*combatlog.Result, bool) {
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	logFile, _, err := r.FormFile("log")

	if err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			log.Printf("upload too large: %v", err)
			http.Error(w, "log file is too large", http.StatusRequestEntityTooLarge)
			return nil, false
		}
		log.Printf("reading uploaded file: %v", err)
		http.Error(w, "missing log file in field \"log\"", http.StatusBadRequest)
		return nil, false
	}

	defer logFile.Close()

	result, err := combatlog.Parse(logFile)
	if err != nil {
		log.Printf("Error parsing log file %v\n", err)
		http.Error(w, "Error parsing log file", http.StatusInternalServerError)
		return nil, false
	}

	return result, true
}
