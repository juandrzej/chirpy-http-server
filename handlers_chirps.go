package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/juandrzej/chirpy-http-server/internal/auth"
	"github.com/juandrzej/chirpy-http-server/internal/database"
)

// Chirp struct to pass all needed data to func responder.
type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

func cleanChirp(body string) string {
	forbiddenWords := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(body, " ")
	for i, word := range words {
		for _, forbiddenWord := range forbiddenWords {
			if strings.ToLower(word) == forbiddenWord {
				words[i] = "****"
			}
		}
	}
	return strings.Join(words, " ")
}

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	// Take parameteres from request and put them into struct.
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameteres")
		return
	}

	// Check if chirp passes internal policy.
	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
	params.Body = cleanChirp(params.Body)

	respondWithJSON(w, http.StatusOK, map[string]string{"cleaned_body": params.Body})
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	// Take parameteres from request and put them into struct.
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameteres")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	params.Body = cleanChirp(params.Body)

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get bearer token from header.")
		return
	}
	tokenID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT.")
		return
	}

	dbChirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   params.Body,
		UserID: tokenID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Chirp could not be created.")
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserId:    dbChirp.UserID,
	})
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps.")
		return
	}

	responseChirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		responseChirps = append(responseChirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserId:    dbChirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, responseChirps)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse chirp ID.")
		return
	}

	dbChirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't find chirp.")
		return
	}

	responseChirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserId:    dbChirp.UserID,
	}

	respondWithJSON(w, http.StatusOK, responseChirp)
}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	// Parse chirp ID.
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse chirp ID.")
		return
	}

	// Validate tokens.
	jwt, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not get refresh token")
		return
	}
	jwtId, err := auth.ValidateJWT(jwt, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not validate JWT")
		return
	}

	// Fetch chirp from DB by its ID.
	dbChirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't find chirp.")
		return
	}

	if dbChirp.UserID != jwtId {
		respondWithError(w, http.StatusForbidden, "Not your chirp doggg.")
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), dbChirp.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not delete chirp.")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
