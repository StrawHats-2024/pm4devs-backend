package user

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"pm4devs-backend/services/auth"
	"pm4devs-backend/types"
	"pm4devs-backend/utils"
	"time"
)

func (s *Handler) handleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	var req types.LoginPayload
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	err = utils.ValidateRequestBody(req, w)
	if err != nil {
		return err
	}
	user, err := s.store.GetUserByEmail(req.Email)
	// TODO: Update error message to when email not found
	if err != nil {
		return err
	}

	if !user.ValidPassword(req.Password) {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.ApiError{Error: "Wrong password or email"})
		return nil
	}
	token, err := auth.CreateJWT(user)
	if err != nil {
		log.Fatal(err)
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Path:    "/",
		Expires: time.Now().Add(time.Hour * 24),
	})
	err = s.store.UpdateLastLogin(user.UserID)
	if err != nil {
		return err
	}
	return utils.WriteJSON(w, http.StatusOK, types.LoginResponse{
		Token:  token,
		UserId: int64(user.UserID),
	})

}
