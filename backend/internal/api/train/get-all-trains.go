package train

import (
	"better-uptime/common/middleware"
	"better-uptime/common/util"
	"fmt"
	"net/http"
)

func (h *Handler) GetAllTrain(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	payload, err := middleware.GetFirebasePayloadFromContext(ctx)
	if err != nil {
		util.ErrorJson(w, err)
		return
	}

	role := payload.Role
	fmt.Print(role)
	if role != "ADMIN" {
		util.ErrorJson(w, util.ErrUnauthorized)
		return
	}

	trainWithSchedule, err := h.store.GetAllTrain(ctx)
	if err != nil {
		util.ErrorJson(w, err)
		return
	}

	response := map[string]interface{}{
		"message": "All trains with schedule",
		"data":    trainWithSchedule,
	}

	util.WriteJson(w, http.StatusOK, response)
}
