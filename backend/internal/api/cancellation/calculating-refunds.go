package cancellation

import (
	"better-uptime/common/middleware"
	"better-uptime/common/util"
	"net/http"
)

func (h *Handler) CalculatingRefundAmount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	payload, err := middleware.GetFirebasePayloadFromContext(ctx)
	if err != nil {
		util.ErrorJson(w, util.ErrUnauthorized)
	}

	userId:= payload.UserId
	
}
