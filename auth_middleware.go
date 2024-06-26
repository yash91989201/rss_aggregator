package main

import (
	"fmt"
	"net/http"

	"github.com/yash91989201/rss_aggregator/internal/auth"
	"github.com/yash91989201/rss_aggregator/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *apiConfig) auth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		apiKey, err := auth.GetApiKey(r.Header)
		if err != nil {
			respondWithError(w, 404, fmt.Sprintf("Unauthorized: %s", err))
			return
		}

		user, err := apiCfg.DB.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, 404, fmt.Sprintf("User not found: %s", err))
			return
		}

		handler(w, r, user)
	}
}
