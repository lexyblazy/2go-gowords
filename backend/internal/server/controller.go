package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	// "log"
	"net/http"
	"strings"

	"github.com/lexyblazy/gowords/internal/helpers"
	"golang.org/x/crypto/bcrypt"
)

type body struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func normalizeUserName(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func (s *Server) login(r *http.Request) (any, int, error) {

	if r.Method != http.MethodPost {
		return nil, http.StatusNotFound, errors.New("not found")
	}

	var reqBody body

	err := json.NewDecoder(r.Body).Decode(&reqBody)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	user, err := s.db.GetUserByUsername(normalizeUserName(reqBody.Username))

	if errors.Is(err, sql.ErrNoRows) {
		return nil, http.StatusUnauthorized, errors.New("account not found")
	}

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqBody.Password))

	if err != nil {
		return nil, http.StatusUnauthorized, errors.New("check username/password combination")
	}

	return user, http.StatusOK, nil

}

func (s *Server) register(r *http.Request) (any, int, error) {

	if r.Method != http.MethodPost {
		return nil, http.StatusNotFound, errors.New("not found")
	}

	var reqBody body

	err := json.NewDecoder(r.Body).Decode(&reqBody)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	user, err := s.db.GetUserByUsername(normalizeUserName(reqBody.Username))

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, http.StatusInternalServerError, err
	}

	if user.ID != "" {
		// user already exists
		return nil, http.StatusUnprocessableEntity, errors.New("username is taken")
	}

	id, err := helpers.NewUUIDV4()

	if err != nil {
		return nil, http.StatusInternalServerError, err

	}
	hash, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	newUser, err := s.db.CreateUser(id, normalizeUserName(reqBody.Username), string(hash), reqBody.Username)

	if err != nil {
		return nil, http.StatusInternalServerError, err

	}

	return newUser, http.StatusOK, nil
}

func (s *Server) logout(r *http.Request) (any, int, error) {
	if r.Method != http.MethodPost {
		return nil, http.StatusNotFound, errors.New("not found")

	}

	return "OK", http.StatusOK, nil
}
