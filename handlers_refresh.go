package main

import (
	"net/http"
	"time"

	"github.com/juandrzej/chirpy-http-server/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	// First find the refresh refreshToken.
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not get refresh token")
		return
	}

	// Locate user in the database using the token.
	dbUser, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not find user in database")
		return
	}

	// Create JWT for user.
	jwt, err := auth.MakeJWT(dbUser.ID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not create JWT")
		return
	}

	// Respond with new token for the user.
	type response struct {
		Token string `json:"token"`
	}
	respondWithJSON(w, http.StatusOK, response{Token: jwt})
}
func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	// First find the refresh refreshToken.
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not get refresh token")
		return
	}

	// Revoke the user refresh token.
	_, err = cfg.db.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not revoke token")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
