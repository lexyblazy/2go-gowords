package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
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

func getCookie(r *http.Request, name string) string {
	cookie, err := r.Cookie(name)
	if err != nil {
		return ""
	}
	return cookie.Value
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

	user, err := s.db.GetUserByUsername(reqBody.Username)

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

	// cache user moniker
	s.rs.CacheUserMoniker(r.Context(), user.ID, user.Moniker)

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

	user, err := s.db.GetUserByUsername(reqBody.Username)

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

	newUser, err := s.db.CreateUser(id, reqBody.Username, string(passHash), reqBody.Username, rcHash)

	if err != nil {
		return nil, http.StatusInternalServerError, err

	}

	err = s.createSession(r, w, newUser)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// cache user moniker
	s.rs.CacheUserMoniker(r.Context(), newUser.ID, newUser.Moniker)

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

	sessionToken := getCookie(r, "session")

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

	user, err := s.db.GetUserByUsername(reqBody.Username)

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

func (s *Server) getLeaderboards(r *http.Request, w http.ResponseWriter) (any, int, error) {

	if r.Method != http.MethodGet {
		return nil, http.StatusNotFound, errors.New("not found")
	}

	daily, err := s.rs.GetDailyLeaderBoard(r.Context())

	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("failed to get daily leaderboard")
	}

	weekly, err := s.rs.GetWeeklyLeaderboard(r.Context())

	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("failed to get weekly leaderboard")
	}

	allTimeHighScores, err := s.rs.GetAllTimeHighScores(r.Context())

	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("failed to get weekly leaderboard")
	}

	res := make(map[string]any)

	res["daily"] = daily
	res["weekly"] = weekly
	res["highScores"] = allTimeHighScores

	return res, http.StatusOK, nil
}
