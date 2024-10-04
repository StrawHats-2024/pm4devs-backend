package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"pm4devs.strawhats/internal/xerrors"
)

// ============================================================================
// Write JSON
// ============================================================================

// Writes JSON to the client
func (rest *Rest) WriteJSON(
	w http.ResponseWriter,
	op string,
	status int,
	data Envelope,
) {
	// Always set the content type as JSON before any other operation
	w.Header().Set("Content-Type", "application/json")

	// Marshal the data to JSON
	response, err := json.Marshal(data)

	// If an error occurs while marshalling, handle it before setting headers
	if err != nil {
		// Log the error using your error handling system
		wrappedError := fmt.Errorf("%w: %v", xerrors.ErrServerInternal, err)
		serverError := xerrors.ServerError(op, wrappedError)
		rest.Logger.Error(serverError.Error())

		// Since there was an error, set the status to 500 (Internal Server Error)
		w.WriteHeader(http.StatusInternalServerError)

		// Return a JSON response for the error, ensuring the content-type is still correct
		w.Write([]byte(`{"error": "Internal server error"}`))
		return
	}

	// Set the response status code and write the JSON response
	w.WriteHeader(status)
	w.Write(response)
}

// ============================================================================
// Read JSON
// ============================================================================

// Reads the request body into the given destination or returns an error
func (rest *Rest) ReadJSON(w http.ResponseWriter, r *http.Request, op string, dst any) *xerrors.AppError {
	r.Body = http.MaxBytesReader(w, r.Body, 1_048_576)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&dst); err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			return xerrors.ClientError(
				http.StatusBadRequest,
				fmt.Sprintf("Request body is malformed (at position %d)", syntaxError.Offset),
				op,
				fmt.Errorf("%w: %v", xerrors.ErrBadRequest, err),
			)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return xerrors.ClientError(
				http.StatusBadRequest,
				"Request body contains badly-formed JSON (at EOF)",
				op,
				fmt.Errorf("%w: %v", xerrors.ErrBadRequest, err),
			)

		case errors.As(err, &unmarshalTypeError):
			return xerrors.ClientError(
				http.StatusBadRequest,
				fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset),
				op,
				fmt.Errorf("%w: %v", xerrors.ErrBadRequest, err),
			)

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return xerrors.ClientError(
				http.StatusBadRequest,
				fmt.Sprintf("Request body contains unknown field %s", fieldName),
				op,
				fmt.Errorf("%w: %v", xerrors.ErrBadRequest, err),
			)

		case errors.Is(err, io.EOF):
			return xerrors.ClientError(
				http.StatusBadRequest,
				"Request body cannot be empty",
				op,
				fmt.Errorf("%w: %v", xerrors.ErrBadRequest, err),
			)

		case err.Error() == "http: request body too large":
			return xerrors.ClientError(
				http.StatusBadRequest,
				"Request body must not be larger than 1MB",
				op,
				fmt.Errorf("%w: %v", xerrors.ErrEntityTooLarge, err),
			)

		default:
			return xerrors.ServerError(
				op,
				fmt.Errorf("%w: %v", xerrors.ErrServerInternal, err),
			)
		}
	}

	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return xerrors.ClientError(
			http.StatusBadRequest,
			"Request body must only contain a single JSON object",
			op,
			fmt.Errorf("%w: %v", xerrors.ErrBadRequest, err),
		)
	}

	return nil
}
