package server

import (
	"errors"
	"net/http"
)

func (s *Server) login(r *http.Request) (any, int, error) {

	if r.Method == http.MethodPost {

		return "OK", http.StatusOK, nil
	}

	return nil, http.StatusNotFound, errors.New("not found")

}

func (s *Server) register(r *http.Request) (any, int, error) {

	if r.Method == http.MethodPost {

		return "OK", http.StatusOK, nil
	}

	return nil, http.StatusNotFound, errors.New("not found")
}

func (s *Server) logout(r *http.Request) (any, int, error) {
	if r.Method == http.MethodPost {

		return "OK", http.StatusOK, nil
	}

	return nil, http.StatusNotFound, errors.New("not found")
}
