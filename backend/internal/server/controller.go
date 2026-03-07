package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/lexyblazy/gowords/internal/helpers"
	"github.com/lexyblazy/gowords/internal/store"
)

type body struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	RecoveryCode string `json:"recoveryCode"`
}

func normalizeUserName(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func (s *Server) createSession(r *http.Request, w http.ResponseWriter, user store.UserEntity) error {
	sessionToken, err := helpers.NewUUIDV4()

	if err != nil {
		return err
	}

	sessionAge := 24 * time.Hour
	err = s.rs.Set(r.Context(), fmt.Sprintf("sessions:%s", sessionToken), user.ID, sessionAge)

	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:     "session",
		Value:    sessionToken,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   int(sessionAge.Seconds()),
	}

	http.SetCookie(w, cookie)

	return nil
}

func (s *Server) login(r *http.Request, w http.ResponseWriter) (any, int, error) {

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

	err = s.createSession(r, w, user)

	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("failed to create session")
	}

	var res map[string]string = make(map[string]string)

	res["id"] = user.ID
	res["moniker"] = user.Moniker
	res["username"] = user.Username

	return res, http.StatusOK, nil

}

func (s *Server) register(r *http.Request, w http.ResponseWriter) (any, int, error) {

	if r.Method != http.MethodPost {
		return nil, http.StatusNotFound, errors.New("not found")
	}

	var reqBody body

	err := json.NewDecoder(r.Body).Decode(&reqBody)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if len(reqBody.Username) < 3 {
		return nil, http.StatusUnprocessableEntity, errors.New("username must be at least 3 characters long")
	}

	if len(reqBody.Password) < 6 {
		return nil, http.StatusUnprocessableEntity, errors.New("password must be at least 6 characters long")
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
	passHash, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	code, _ := helpers.GenerateRecoveryCode()
	rcHash := helpers.HashRecoveryCode(code)

	newUser, err := s.db.CreateUser(id, normalizeUserName(reqBody.Username), string(passHash), reqBody.Username, rcHash)

	if err != nil {
		return nil, http.StatusInternalServerError, err

	}

	err = s.createSession(r, w, newUser)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	var res map[string]string = make(map[string]string)

	res["id"] = newUser.ID
	res["moniker"] = newUser.Moniker
	res["username"] = newUser.Username
	res["recoveryCode"] = code

	return res, http.StatusOK, nil
}

func (s *Server) logout(r *http.Request, w http.ResponseWriter) (any, int, error) {
	if r.Method != http.MethodPost {
		return nil, http.StatusNotFound, errors.New("not found")

	}

	cookie, err := r.Cookie("session")

	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("something went wrong")
	}

	sessionToken := cookie.Value

	if sessionToken == "" {
		return nil, http.StatusUnauthorized, errors.New("no session found")
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	s.rs.Delete(r.Context(), fmt.Sprintf("sessions:%s", sessionToken))

	return "OK", http.StatusOK, nil
}

func (s *Server) resetPassword(r *http.Request, w http.ResponseWriter) (any, int, error) {
	if r.Method != http.MethodPost {
		return nil, http.StatusNotFound, errors.New("not found")

	}

	var reqBody body

	err := json.NewDecoder(r.Body).Decode(&reqBody)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	user, err := s.db.GetUserByUsername(normalizeUserName(reqBody.Username))

	if err != nil || errors.Is(err, sql.ErrNoRows) {
		return nil, http.StatusUnauthorized, errors.New("account not found")
	}

	if len(reqBody.RecoveryCode) < 1 {
		return nil, http.StatusUnauthorized, errors.New("recovery code is required")
	}

	if !helpers.ValidateRecoveryCode(reqBody.RecoveryCode, user.RecoveryHash) {
		return nil, http.StatusUnauthorized, errors.New("recovery code is incorrect")
	}

	if len(reqBody.Password) < 6 {
		return nil, http.StatusUnprocessableEntity, errors.New("password must be at least 6 characters long")
	}

	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), bcrypt.DefaultCost)

	err = s.db.UpdatePassword(user.ID, string(newPasswordHash))

	if err != nil {
		return nil, http.StatusUnauthorized, errors.New("failed to update password")
	}

	err = s.createSession(r, w, user)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	var res map[string]string = make(map[string]string)

	res["id"] = user.ID
	res["moniker"] = user.Moniker
	res["username"] = user.Username

	return res, http.StatusOK, nil

}
