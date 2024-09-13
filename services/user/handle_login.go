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

	"github.com/go-playground/validator/v10"
)

func (s *Handler) handleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	var req types.LoginPayload
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println("err: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return err
	}
	if err := utils.Validate.Struct(req); err != nil {
		errors := err.(validator.ValidationErrors)
		return utils.WriteJSON(w, http.StatusBadRequest,
			utils.ApiError{Error: fmt.Errorf("invalid payload: %v", errors).Error()})
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
	err = utils.WriteJSON(w, http.StatusOK, types.LoginResponse{
		Token:  token,
		UserId: int64(user.UserID),
	})
	return nil

}
