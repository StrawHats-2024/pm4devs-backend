package user

import (
	"encoding/json"
	"net/http"
	"pm4devs-backend/services/auth"
	"pm4devs-backend/utils"
)

func (s *Handler) handleVerifyToken(w http.ResponseWriter, r *http.Request) error {

	var req struct {
		Token string `json:"token" validate:"required"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}
	err = utils.ValidateRequestBody(req, w)
	if err != nil {
		return err
	}
	user, err := auth.ValidateToken(req.Token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return err
	}
	type TokenVerificationResponse struct {
		Valid   bool   `json:"valid"`
		UserID  int    `json:"user_id"`
		Message string `json:"message"`
	}
	return utils.WriteJSON(w, http.StatusOK, TokenVerificationResponse{
		Valid:   true,
		UserID:  user.UserId,
		Message: "Token is valid",
	})

}
