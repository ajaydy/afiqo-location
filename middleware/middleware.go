package middleware

import (
	"afiqo-location/helpers"
	"afiqo-location/session"
	"context"
	"github.com/gomodule/redigo/redis"
	"log"
	"net/http"
)

func SessionMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		sessionID := r.Header.Get("session")

		session := session.Session{
			SessionKey: sessionID,
		}

		sessionData, err := session.Get(ctx)

		if err != nil {
			if err == redis.ErrNil {
				helpers.ErrorResponse(w, helpers.UnauthorizedMessage, http.StatusUnauthorized)
				return
			}
			helpers.ErrorResponse(w, helpers.InternalServerError, http.StatusInternalServerError)
			return

		}

		ctx = context.WithValue(ctx, "user_id", sessionData.UserID.String())
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func RolesMiddleware(next http.Handler, roles ...string) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		sessionID := r.Header.Get("session")

		session := session.Session{
			SessionKey: sessionID,
		}

		sessionData, err := session.Get(ctx)

		if err != nil {
			if err == redis.ErrNil {
				helpers.ErrorResponse(w, helpers.UnauthorizedMessage, http.StatusUnauthorized)
				return
			}
			helpers.ErrorResponse(w, helpers.InternalServerError, http.StatusInternalServerError)
			return

		}

		validRole := false

		for _, role := range roles {
			if role == sessionData.Role {
				validRole = true
				break
			}
		}

		if !validRole {
			helpers.ErrorResponse(w, helpers.ForbiddenMessage, http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
func LoggingMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		reqURI := r.RequestURI
		if reqURI == "" {
			reqURI = r.URL.RequestURI()
		}

		log.Println(reqURI)

		ctx := r.Context()
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)

	})

}
