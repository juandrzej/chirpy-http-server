package main

import (
	"net/http"
	"time"

	"github.com/juandrzej/chirpy-http-server/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	jwt, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not find JWT in headers.")
		return
	}

	dbUser, err := cfg.db.GetUserFromRefreshToken(r.Context(), jwt)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not find user in database.")
		return
	}

	jwtId, err := auth.MakeJWT(dbUser.ID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not create JWT.")
		return
	}

	type response struct {
		Token string `json:"token"`
	}
	respondWithJSON(w, http.StatusOK, response{Token: jwtId})
}
func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	jwt, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not find JWT in headers.")
		return
	}

	_, err = cfg.db.RevokeRefreshToken(r.Context(), jwt)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not revoke token.")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
