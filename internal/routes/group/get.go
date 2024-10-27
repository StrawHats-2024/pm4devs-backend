package group

import (
	"net/http"

	"pm4devs.strawhats/internal/rest"
	"pm4devs.strawhats/internal/routes/middleware"
)

const ListUserGroupRoute = "/v1/groups/user"

func (app *Group) listUserGroups(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		app.rest.MethodNotAllowed(w, r, "GET")
		return
	}

	currUser := middleware.ContextGetUser(r)
	groups, err := app.group.GetGroupsByUserID(currUser.ID)
	if err != nil {
		app.rest.Error(w, err)
		return
	}
	app.rest.WriteJSON(w, "group.createNew", http.StatusOK, rest.Envelope{
		"Message": "Success!",
		"data":    groups,
	})
}
