package group

import (
	"net/http"

	"pm4devs.strawhats/internal/rest"
	"pm4devs.strawhats/internal/routes/middleware"
	"pm4devs.strawhats/internal/validator"
)

const AddUserToGroupRoute = "/v1/groups/add_user"
const RemoveUserFromGroupRoute = "/v1/groups/remove_user"

func (app *Group) addUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.rest.MethodNotAllowed(w, r, "POST")
		return
	}

	var input struct {
		GroupName string `json:"group_name"`
		UserEmail string `json:"user_email"`
	}

	// Parse request
	if err := app.rest.ReadJSON(w, r, "group.addUser", &input); err != nil {
		app.rest.Error(w, err)
		return
	}

	// Validate parameters
	v := validator.New()
	v.Check(len(input.GroupName) > 0, "group_name", "must be provided")
	v.Check(len(input.UserEmail) > 0, "user_email", "must be provided")
	if err := v.Valid("group.addUser"); err != nil {
		app.rest.Error(w, err)
		return
	}

	currUser := middleware.ContextGetUser(r)
	user, err := app.users.GetByEmail(input.UserEmail)
	if err != nil {
		app.rest.Error(w, err)
		return
	}
	group, err := app.group.GetByGroupName(input.GroupName)
	if err != nil {
		app.rest.Error(w, err)
		return
	}
	currGroup, err := app.group.GetByGroupID(group.ID)
	if err != nil {
		app.rest.Error(w, err)
		return
	}
  // check if the user is creator of the group
	if currUser.ID != currGroup.CreatorID {
		app.rest.WriteJSON(w, "group.addUser", http.StatusUnauthorized, rest.Envelope{
			"Message": "Only owner can add member to the group",
		})
		return
	}
	err = app.group.AddUser(group.ID, user.ID)
	if err != nil {
		app.rest.Error(w, err)
		return
	}
	app.rest.WriteJSON(w, "group.addUser", http.StatusOK, rest.Envelope{
		"Message": "Success!",
	})

}

func (app *Group) removeUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		app.rest.MethodNotAllowed(w, r, "POST")
		return
	}

	var input struct {
		GroupName string `json:"group_name"`
		UserEmail string `json:"user_email"`
	}

	// Parse request
	if err := app.rest.ReadJSON(w, r, "group.removeUser", &input); err != nil {
		app.rest.Error(w, err)
		return
	}

	// Validate parameters
	v := validator.New()
	v.Check(len(input.GroupName) > 0, "group_name", "must be provided")
	v.Check(len(input.UserEmail) > 0, "user_email", "must be provided")
	if err := v.Valid("group.removeUser"); err != nil {
		app.rest.Error(w, err)
		return
	}

	currUser := middleware.ContextGetUser(r)
	user, err := app.users.GetByEmail(input.UserEmail)
	if err != nil {
		app.rest.Error(w, err)
		return
	}
	currGroup, err := app.group.GetByGroupName(input.GroupName)
	if err != nil {
		app.rest.Error(w, err)
		return
	}
	if currUser.ID != currGroup.CreatorID {
		app.rest.WriteJSON(w, "group.removeUser", http.StatusUnauthorized, rest.Envelope{
			"Message": "Only owner can remove member to the group",
		})
		return
	}
	if user.ID == currGroup.CreatorID {
		app.rest.WriteJSON(w, "group.removeUser", http.StatusBadRequest, rest.Envelope{
			"Message": "Invalid user_id, trying to make create a member",
		})
		return
	}
	err = app.group.RemoveUser(currGroup.ID, user.ID)
	if err != nil {
		app.rest.Error(w, err)
		return
	}
	app.rest.WriteJSON(w, "group.removeUser", http.StatusOK, rest.Envelope{
		"Message": "Success!",
	})

}
