package network

import (
	"encoding/json"
	"log"
	"net/http"
)

type Responses struct{}

// main response function
func res(w http.ResponseWriter, _ *http.Request, statusCode int, message string, data any) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")

	resp := make(map[string]any)

	if statusCode >= 200 && statusCode < 400 {
		resp["success"] = true
	} else {
		resp["success"] = false
	}

	resp["data"] = data

	resp["message"] = message

	jsonRes, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in json err: %s", err)
	}

	w.Write(jsonRes)
}

/**
sucessful
*/

// get request
func (Responses) RespondOk(w http.ResponseWriter, r *http.Request, data any, message ...string) {
	mes := "successful"
	if len(message) > 0 {
		mes = message[0]
	}
	res(w, r, http.StatusOK, mes, data)
}

// post request
func (Responses) RespondCreated(w http.ResponseWriter, r *http.Request, data any, message ...string) {
	msg := "created successfully"
	if len(message) > 0 {
		msg = message[0]
	}
	res(w, r, http.StatusCreated, msg, data)
}

// accepted request
func (Responses) RespondAccepted(w http.ResponseWriter, r *http.Request) {
	res(w, r, http.StatusAccepted, "process initiated", nil)
}

// put request
func (Responses) RespondUpdated(w http.ResponseWriter, r *http.Request) {
	res(w, r, http.StatusNoContent, "updated", nil)
}

// delete request
func (Responses) RespondDeleted(w http.ResponseWriter, r *http.Request) {
	res(w, r, http.StatusNoContent, "deleted", nil)
}

/**
end sucessful
*/

/**
failed
*/

// bad request
func (Responses) RepondBadRequest(w http.ResponseWriter, r *http.Request, message...string) {
	msg := "bad request"
	if len(message) > 0 {
		msg = message[0]
	}
	res(w, r, http.StatusBadRequest, msg, nil)
}

// unauthorized request
func (Responses) RepondUnauthorized(w http.ResponseWriter, r *http.Request) {
	res(w, r, http.StatusUnauthorized, "unauthorized", nil)
}

// forbidden
func (Responses) RepondForbidden(w http.ResponseWriter, r *http.Request) {
	res(w, r, http.StatusForbidden, "Forbidden", nil)
}

/**
end failed
*/
