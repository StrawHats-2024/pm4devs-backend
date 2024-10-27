package group

import (
	"net/http"
	"testing"

	"pm4devs.strawhats/internal/app"
	"pm4devs.strawhats/internal/assert"
	"pm4devs.strawhats/internal/mocks"
	"pm4devs.strawhats/internal/routes/group"
	"pm4devs.strawhats/internal/routes/middleware"
	"pm4devs.strawhats/internal/routes/utils"
)

func TestCRUDGroup(t *testing.T) {
	assert.Integration(t)
	app := mocks.App(t)
	handler := groupHandler(app)
	authHandler := utils.AuthHandler(app)

	credentials := `{"email": "test@example.com", "password": "password"}`
	credentialsTwo := `{"email": "test2@example.com", "password": "password"}`

	assert.Check(t, utils.RegisterUser(authHandler, credentials))
	token := utils.LoginUser(authHandler, credentials)

	assert.Check(t, utils.RegisterUser(authHandler, credentialsTwo))
	tokenTwo := utils.LoginUser(authHandler, credentialsTwo)

	assert.Check(t, len(token) > 0)

	type responseMessage struct {
		Error   map[string]string `json:"error"`
		Message string            `json:"message"`
		Data    map[string]any    `json:"data"`
	}

	tests := []assert.HandlerTestCase[responseMessage]{
		{
			Name:   "AuthRequired",
			Status: http.StatusUnauthorized,
			Method: http.MethodGet,
		},
		{
			Name:   "MethodNotAllowed/PUT",
			Auth:   token,
			Status: http.StatusMethodNotAllowed,
			Method: http.MethodPut,
		},
		{
			Name:   "InvalidBody/POST",
			Auth:   token,
			Status: http.StatusBadRequest,
			Method: http.MethodPost,
		},
		{
			Name:   "InvalidGroupNameLength/POST",
			Auth:   token,
			Status: http.StatusUnprocessableEntity,
			Body:   `{"group_name": "tes"}`,
			Method: http.MethodPost,
			FN: func(t *testing.T, result responseMessage) {
				// t.Errorf("%+v", result)
				assert.Equal(t, result.Error["group_name"], "must be provided & at least of 5 charators long")
			},
		},
		{
			Name:   "ValidBody/POST",
			Auth:   token,
			Status: http.StatusCreated,
			Body:   `{"group_name": "testgroup"}`,
			Method: http.MethodPost,
			FN: func(t *testing.T, result responseMessage) {
				assert.Equal(t, result.Message, "Success!")
				assert.NotEqual(t, result.Data["group_name"], "")
			},
		},
		{
			Name:   "DuplicateGroupName/POST",
			Auth:   token,
			Status: http.StatusConflict,
			Body:   `{"group_name": "testgroup"}`,
			Method: http.MethodPost,
		},
		{
			Name:   "ValidBody2/POST",
			Auth:   token,
			Status: http.StatusCreated,
			Body:   `{"group_name": "testgroup2"}`,
			Method: http.MethodPost,
		},
		// Delete group by id
		{
			Name:   "EmptyBody/DEL",
			Auth:   token,
			Status: http.StatusBadRequest,
			Method: http.MethodDelete,
		},
		{
			Name:   "InvalidBody/DEL",
			Auth:   token,
			Body:   `{"group_name": }`,
			Status: http.StatusBadRequest,
			Method: http.MethodDelete,
		},
		{
			Name:   "ValidBodyNotOwner/DEL",
			Auth:   tokenTwo,
			Body:   `{"group_name": "testgroup2"}`,
			Status: http.StatusUnauthorized,
			Method: http.MethodDelete,
		},
		{
			Name:   "ValidBody/DEL",
			Auth:   token,
			Body:   `{"group_name": "testgroup2"}`,
			Status: http.StatusNoContent,
			Method: http.MethodDelete,
		},
		{
			Name:   "ValidBodyInvalidGroupID/DEL",
			Auth:   token,
			Body:   `{"group_name": "testgroup2"}`,
			Status: http.StatusNotFound,
			Method: http.MethodDelete,
		},

		// update group name

		{
			Name:   "EmptyBody/UPDATE",
			Auth:   token,
			Status: http.StatusBadRequest,
			Method: http.MethodPatch,
		},
		{
			Name:   "InvalidGroupName/UPDATE",
			Auth:   tokenTwo,
			Body:   `{"group_name": "testgrou", "new_group_name": "update"}`,
			Status: http.StatusNotFound,
			Method: http.MethodPatch,
		},
		{
			Name:   "ValidBodyNotOwner/UPDATE",
			Auth:   tokenTwo,
			Body:   `{"group_name": "testgroup", "new_group_name": "newgroup"}`,
			Status: http.StatusUnauthorized,
			Method: http.MethodPatch,
		},
		{
			Name:   "ValidBodyOwner/UPDATE",
			Auth:   token,
			Body:   `{"group_name": "testgroup", "new_group_name": "newgroup"}`,
			Status: http.StatusOK,
			Method: http.MethodPatch,
		},

		// Get group by id
		{
			Name:   "EmptyBody/GET",
			Auth:   token,
			Status: http.StatusBadRequest,
			Method: http.MethodGet,
		},
		{
			Name:   "InvalidGroupID/GET",
			Auth:   tokenTwo,
			Body:   `{"group_name": "ggg"}`,
			Status: http.StatusNotFound,
			Method: http.MethodGet,
		},
		{
			Name:   "ValidBody/GET",
			Auth:   token,
			Body:   `{"group_name": "newgroup"}`,
			Status: http.StatusOK,
			Method: http.MethodGet,
		},
	}

	for _, tc := range tests {
		assert.RunHandlerTestCase(t, handler, tc.Method, group.CRUDGroupRoute, tc)
	}

}

func groupHandler(app *app.App) http.HandlerFunc {
	handler := func() http.Handler {
		mux := http.NewServeMux()

		middleware := middleware.New(app)
		secrets := group.New(app)
		secrets.Route(mux, middleware)

		return middleware.User(mux)
	}()

	return func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	}
}
