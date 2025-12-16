package middleware

import (
	"better-uptime/common/firebase"
	"better-uptime/common/util"
	db "better-uptime/internal/db/sqlc"
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sirupsen/logrus"
)

type tokenPayloadKeyType string

const TokenPayloadKey tokenPayloadKeyType = "auth-payload"

// TokenMiddleware verifies Firebase token, upserts user, and sets payload in context
func TokenMiddleware(store db.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				util.ErrorJson(w, util.ErrInvalidToken)
				logrus.Error("Authorization header missing")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				util.ErrorJson(w, util.ErrInvalidToken)
				logrus.Error("Invalid Authorization header format")
				return
			}

			idToken := parts[1]

			ctx := r.Context()
			var payload firebase.FirebasePayload

			// Dev mode: frontend bypass
			if idToken == "frontend" {
				payload = firebase.FirebasePayload{
					Email:  "test@example.com",
					UID:    "0b927d97-782a-4c82-b9d2-e4e06774ed37",
					UserId: uuid.MustParse("0b927d97-782a-4c82-b9d2-e4e06774ed37"),
					Role:   "ADMIN",
				}
				logrus.Info("Frontend bypass: using test user payload")
				ctx = context.WithValue(ctx, TokenPayloadKey, payload)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// Verify Firebase token
			fbPayload, err := firebase.VerifyFirebaseIDToken(ctx, idToken)
			if err != nil {
				util.ErrorJson(w, util.ErrInvalidToken)
				logrus.WithError(err).Error("Failed to verify Firebase ID token")
				return
			}

			logrus.Infof("Firebase token verified. Email: %s, UID: %s", fbPayload.Email, fbPayload.UID)

			// Set defaults for provider and fullname
			if fbPayload.Provider == "" {
				fbPayload.Provider = "password"
			}
			if fbPayload.Fullname == "" {
				if at := strings.Index(fbPayload.Email, "@"); at > 0 {
					fbPayload.Fullname = fbPayload.Email[:at]
				} else {
					fbPayload.Fullname = "User"
				}
			}

			// Upsert user in DB
			user, err := store.FindOrCreateUser(ctx, db.FindOrCreateUserParams{
				Email:        fbPayload.Email,
				Fullname:     fbPayload.Fullname,
				Provider:     db.Provider(fbPayload.Provider),
				Phone:        pgxText(fbPayload.Phone),
				PasswordHash: pgxText(""), // Firebase doesn't use password here
			})
			if err != nil {
				util.ErrorJson(w, err)
				logrus.WithError(err).Error("Failed to find or create user")
				return
			}

			// Attach internal UserId to payload
			fbPayload.UserId = user.ID
			if fbPayload.Role == "" {
				fbPayload.Role = "user"
			}

			logrus.Infof("User found/created: %v", user.ID)

			// Set payload in request context
			ctx = context.WithValue(ctx, TokenPayloadKey, fbPayload)

			// Proceed to next handler
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Helper to convert string to pgtype.Text
func pgxText(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: s != ""}
}

// Extract FirebasePayload from context
func GetFirebasePayloadFromContext(ctx context.Context) (firebase.FirebasePayload, error) {
	raw := ctx.Value(TokenPayloadKey)
	if raw == nil {
		return firebase.FirebasePayload{}, errors.New("missing auth token payload in context")
	}

	payload, ok := raw.(firebase.FirebasePayload)
	if !ok {
		return firebase.FirebasePayload{}, errors.New("invalid auth token payload type")
	}

	return payload, nil
}
