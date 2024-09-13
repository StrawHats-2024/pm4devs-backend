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

func (s *Handler) handleRegister(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	var req types.RegisterUserPayload
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

	user, err := utils.NewUser(req.Email, req.Username, req.Password)
	if err != nil {
		return err
	}
	userId, err := s.store.CreateUser(user)
	if err != nil {
		fmt.Printf("user: %+v", user)
		fmt.Println("err: ", err)
		return err
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
	err = utils.WriteJSON(w, http.StatusOK, types.RegisterUserResponse{
		Token:   token,
		Message: "User Registered successfully.",
		UserId:  int64(userId),
	})
	return nil

}
