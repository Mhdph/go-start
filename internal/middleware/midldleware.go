package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/mhdph/go-start/internal/store"
	"github.com/mhdph/go-start/internal/store/tokens"
	"github.com/mhdph/go-start/internal/utils"
)

type UserMiddlware struct {
	UserStore store.UserStore
}

type contextKey string

const userContextKey = contextKey("user")

func SetUser(r *http.Request, user *store.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)

}

func GetUser(r *http.Request) *store.User {
	user, ok := r.Context().Value(userContextKey).(*store.User)

	if !ok {
		panic("mising user in request")
	}

	return user

}

func (um *UserMiddlware) Autheniticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("very", "Authorization")
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			r = SetUser(r, store.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"eroor": "invalid authorization header"})
			return
		}

		token := headerParts[1]

		user, err := um.UserStore.GetUserToken(tokens.ScopeAuthentication, token)

		if err != nil {
			utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid token"})
			return

		}
		if user == nil {
			utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid token"})

			return

		}

		r = SetUser(r, user)
		next.ServeHTTP(w, r)
		return

	})
}

func (um *UserMiddlware) RequireUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r)
		if user == nil {
			utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"error": "unauthorized"})
			return
		}
		next.ServeHTTP(w, r)
	})
}
