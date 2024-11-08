package group

import (
	"net/http"

	"pm4devs.strawhats/internal/rest"
	"pm4devs.strawhats/internal/routes/middleware"
	"pm4devs.strawhats/internal/validator"
)

const CRUDGroupRoute = "/v1/groups"

func (app *Group) CRUDRoute(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.get(w, r)

	case http.MethodPost:
		app.createNew(w, r)

	case http.MethodPatch:
		app.update(w, r)

	case http.MethodDelete:
		app.delete(w, r)

	default:
		app.rest.MethodNotAllowed(w, r, "GET, POST, PATCH, DELETE")
	}
}

func (app *Group) createNew(w http.ResponseWriter, r *http.Request) {

	var input struct {
		GroupName string `json:"group_name"`
	}

	// Parse request
	if err := app.rest.ReadJSON(w, r, "group.createNew", &input); err != nil {
		app.rest.Error(w, err)
		return
	}

	// Validate parameters
	v := validator.New()
	v.Check(len(input.GroupName) > 4, "group_name", "must be provided & at least of 5 charators long")
	if err := v.Valid("group.createNew"); err != nil {
		app.rest.Error(w, err)
		return
	}
	currUser := middleware.ContextGetUser(r)
	newGroup, err := app.group.NewRecord(input.GroupName, currUser.ID)
	if err != nil {
		app.rest.Error(w, err)
		return
	}
	app.rest.WriteJSON(w, "group.createNew", http.StatusCreated, rest.Envelope{
		"Message": "Success!",
		"data":    newGroup,
	})
}

func (app *Group) delete(w http.ResponseWriter, r *http.Request) {

	var input struct {
		GroupName string `json:"group_name"`
	}

	// Parse request
	if err := app.rest.ReadJSON(w, r, "group.delete", &input); err != nil {
		app.rest.Error(w, err)
		return
	}

	// Validate parameters
	v := validator.New()
	v.Check(len(input.GroupName) > 0, "group_name", "must be provided")
	if err := v.Valid("group.delete"); err != nil {
		app.rest.Error(w, err)
		return
	}
	currUser := middleware.ContextGetUser(r)
	currGroup, err := app.group.GetGroupUsers(input.GroupName)
	if err != nil {
		app.rest.Error(w, err)
		return
	}
	if currUser.ID != currGroup.CreatorID {
		app.rest.WriteJSON(w, "group.delete", http.StatusUnauthorized, rest.Envelope{
			"Message": "Only owner can delete the group.",
		})
		return
	}
	err = app.group.DeleteByGroupID(currGroup.ID)
	if err != nil {
		app.rest.Error(w, err)
		return
	}
	app.rest.WriteJSON(w, "group.delete", http.StatusNoContent, rest.Envelope{
		"Message": "Success!",
	})
}
func (app *Group) update(w http.ResponseWriter, r *http.Request) {
	var input struct {
		NewGroupName string `json:"new_group_name"`
		GroupName    string `json:"group_name"`
	}

	// Parse request
	if err := app.rest.ReadJSON(w, r, "group.update", &input); err != nil {
		app.rest.Error(w, err)
		return
	}

	// Validate parameters
	v := validator.New()
	v.Check(len(input.NewGroupName) > 4, "new_group_name", "must be provided & at least of 5 charators long")
	v.Check(len(input.GroupName) > 0, "group_name", "must be provided")
	if err := v.Valid("group.update"); err != nil {
		app.rest.Error(w, err)
		return
	}
	currUser := middleware.ContextGetUser(r)
	currGroup, err := app.group.GetGroupUsers(input.GroupName)
	if err != nil {
		app.rest.Error(w, err)
		return
	}
	if currUser.ID != currGroup.CreatorID {
		app.rest.WriteJSON(w, "group.delete", http.StatusUnauthorized, rest.Envelope{
			"Message": "Only owner can delete the group.",
		})
		return
	}
	_, err = app.group.UpdateGroupName(input.NewGroupName, input.GroupName)
	if err != nil {
		app.rest.Error(w, err)
		return
	}
	app.rest.WriteJSON(w, "group.delete", http.StatusOK, rest.Envelope{
		"Message": "Success!",
	})
}
func (app *Group) get(w http.ResponseWriter, r *http.Request) {

	var input struct {
		GroupName string `json:"group_name"`
	}
	// Parse request
	if err := app.rest.ReadJSON(w, r, "group.get", &input); err != nil {
		app.rest.Error(w, err)
		return
	}
	v := validator.New()
	v.Check(len(input.GroupName) > 0, "group_name", "must be provided")
	if err := v.Valid("group.update"); err != nil {
		app.rest.Error(w, err)
		return
	}
	usersInGroup, err := app.group.GetGroupUsers(input.GroupName)
	secretsInGroup, err := app.group.GetGroupSharedSecrets(input.GroupName)
	if err != nil {
		app.rest.Error(w, err)
		return
	}
	app.rest.WriteJSON(w, "group.get", http.StatusOK, rest.Envelope{
		"message": "Success!",
		"data": rest.Envelope{
			"group_id":   usersInGroup.ID,
			"group_name": usersInGroup.Name,
			"created_at": usersInGroup.CreatedAt,
			"creator_id": usersInGroup.CreatorID,
			"users":      usersInGroup.Users,
			"secrets":    secretsInGroup.Secrets,
		},
	})
}
