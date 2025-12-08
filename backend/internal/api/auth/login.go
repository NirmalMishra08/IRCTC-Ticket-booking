package auth

import (
	"net/http"
	"strings"

	"better-uptime/common/firebase"
	"better-uptime/common/logger"
	"better-uptime/common/util"
	db "better-uptime/internal/db/sqlc"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		util.ErrorJson(w, util.ErrInvalidToken)
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		util.ErrorJson(w, util.ErrInvalidToken)
		return
	}
	idToken := parts[1]

	ctx := r.Context()
	var payload firebase.FirebasePayload
	var err error

	// For frontend testing bypass
	if idToken == "frontend" {
		payload = firebase.FirebasePayload{
			Email:    "test@example.com",
			UserId:   uuid.MustParse("0b927d97-782a-4c82-b9d2-e4e06774ed37"),
			UID:      "0b927d97-782a-4c82-b9d2-e4e06774ed37",
			Phone:    "9876543210",
			Provider: "password",
		}
	} else {
		payload, err = firebase.VerifyFirebaseIDToken(ctx, idToken)
		if err != nil {
			util.ErrorJson(w, err)
			return
		}
	}

	var req authRequest
	if err := util.ReadJsonAndValidate(w, r, &req); err != nil {
		util.ErrorJson(w, err)
		return
	}

	// Handle password hash
	var hashedPassword string
	if req.Password != "" {
		pw, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			util.ErrorJson(w, err)
			return
		}
		hashedPassword = string(pw)
	}

	// Upsert user (update every time with latest data)
	phone := payload.Phone
	if req.Phone != "" {
		phone = req.Phone
	}

	provider := req.Provider
	if provider == "" {
		provider = payload.Provider // fallback to payload's provider
	}

	user, err := h.store.FindOrCreateUser(ctx, db.FindOrCreateUserParams{
		Email:        payload.Email,
		Provider:     db.Provider(provider),
		Phone:        pgtype.Text{String: phone, Valid: phone != ""},
		Fullname:     req.FullName,
		PasswordHash: pgtype.Text{String: hashedPassword, Valid: hashedPassword != ""},
	})

	if err != nil {
		logger.Error("FindOrCreateUser failed for email %s: %v", payload.Email, err)
		util.ErrorJson(w, err)
		return
	}

	// Send welcome notification for new users
	// if !userExistedBefore {
	// 	_, _ = h.store.InsertNotification(ctx, db.InsertNotificationParams{
	// 		UserID:           user.ID,
	// 		NotificationType: db.NotificationTypePROFILECREATEDSUCCESSFULLY,
	// 		Category:         db.NotificationTypePROFILECREATEDSUCCESSFULLY,
	// 		EntityID:         user.ID.String(),
	// 		Title:            "Profile created successfully",
	// 		Message:          "Welcome aboard! Your fitness journey starts here.",
	// 	})
	// 	_ = h.firestore.NCPushNotification(
	// 		ctx,
	// 		user.ID.String(),
	// 		"Profile created successfully",
	// 		"Welcome aboard! Your fitness journey starts here.",
	// 		string(db.NotificationTypePROFILECREATEDSUCCESSFULLY),
	// 		user.ID.String(),
	// 		time.Now(),
	// 	)
	// }

	// Respond with user info
	util.WriteJson(w, http.StatusOK, authResponse{
		UserID:   user.ID.String(),
		Provider: string(user.Provider),
		Email:    user.Email,
		FullName: user.Fullname,
	})
}
