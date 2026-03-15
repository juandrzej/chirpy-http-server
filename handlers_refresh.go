package main

import (
	"net/http"

	"github.com/juandrzej/chirpy-http-server/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not get refresh token")
		return

	}

}

