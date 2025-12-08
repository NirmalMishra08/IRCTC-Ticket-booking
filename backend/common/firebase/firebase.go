package firebase

import (
	"context"
	"errors"
	"sync"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

var (
	client   *auth.Client
	initOnce sync.Once
	initErr  error
)

// FirebasePayload holds essential auth user info.
type FirebasePayload struct {
	UID      string    // Firebase UID
	Email    string    // Firebase Email
	Fullname string    // Firebase Full Name
	Phone    string    // Firebase Phone Number
	UserId   uuid.UUID // Your internal DB UUID (populated later)
	Provider string    // Firebase sign-in provider (e.g., password, google, phone)
	Role     string    // Your internal DB role
}

// Initialize Firebase app and client once
func InitFirebaseAuth(serviceAccountPath string) error {
	initOnce.Do(func() {
		opt := option.WithCredentialsFile(serviceAccountPath)
		app, err := firebase.NewApp(context.Background(), nil, opt)
		if err != nil {
			initErr = err
			return
		}
		client, initErr = app.Auth(context.Background())
	})
	return initErr
}

// VerifyFirebaseIDToken verifies the ID token and extracts relevant user info
func VerifyFirebaseIDToken(ctx context.Context, idToken string) (FirebasePayload, error) {
	if client == nil {
		return FirebasePayload{}, errors.New("firebase auth client not initialized")
	}

	token, err := client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return FirebasePayload{}, err
	}

	if token.Expires < time.Now().Unix() {
		return FirebasePayload{}, errors.New("token expired")
	}

	claims := token.Claims

	return FirebasePayload{
		UID:      token.UID,
		Email:    getClaimString(claims, "email"),
		Fullname: getClaimString(claims, "name"),
		Phone:    getClaimString(claims, "phone_number"),
		Provider: getSignInProvider(claims),
	}, nil
}

// getClaimString safely extracts a string from token claims
func getClaimString(claims map[string]interface{}, key string) string {
	if val, ok := claims[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// getSignInProvider extracts the nested `firebase.sign_in_provider` claim
func getSignInProvider(claims map[string]interface{}) string {
	if firebaseClaim, ok := claims["firebase"].(map[string]interface{}); ok {
		if provider, ok := firebaseClaim["sign_in_provider"].(string); ok {
			return provider
		}
	}
	return ""
}
