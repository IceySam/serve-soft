package network

import (
	"encoding/json"
	"log"
	"net/http"
)

type Responses struct {}

// main function
func res(w http.ResponseWriter, r *http.Request, statusCode int, message string, data ...map[string]any)  {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")

	resp := make(map[string]any)
	resp["message"] = message
	if len(data) > 0 {
		resp["data"] = data[0]
	}

	jsonRes, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in json err: %s", err)
	}
	w.Write(jsonRes)
	return
}


/**
	sucessful
*/

// get request
func (Responses) respondOk(w http.ResponseWriter, r *http.Request, message ...string)  {
	mes := "successful"
	if len(message) > 0 {
		mes = message[0]
	}
	res(w,r,http.StatusOK,mes, nil)
	return
}

// post request
func (Responses) respondCreated(w http.ResponseWriter, r *http.Request, result ...map[string]any)  {
	mes := "created successfully"
	data := make(map[string]any)
	if len(result) > 0 {
		data = result[0]
	}
	res(w,r,http.StatusCreated,mes, data)
	return
}

// accepted request
func (Responses) respondAccepted(w http.ResponseWriter, r *http.Request)  {
	res(w,r,http.StatusAccepted,"process initiated", nil)
	return
}

// put request
func (Responses) respondUpdated(w http.ResponseWriter, r *http.Request)  {
	res(w,r,http.StatusNoContent,"updated", nil)
	return
}

// delete request
func (Responses) respondDeleted(w http.ResponseWriter, r *http.Request)  {
	res(w,r,http.StatusNoContent,"deleted", nil)
	return
}

/**
	end sucessful
*/

/**
	failed
*/

// not found
func (Responses) repondBadRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	return
}

// unauthorized request
func (Responses) repondUnauthorized(w http.ResponseWriter, r *http.Request)  {
	res(w,r,http.StatusUnauthorized,"deleted", nil)
	return
}

// forbidden
func (Responses) repondForbidden(w http.ResponseWriter, r *http.Request)  {
	res(w,r,http.StatusForbidden,"Forbidden", nil)
	return
}

/**
	end failed
*/