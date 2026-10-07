package api

import "net/http"

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}
