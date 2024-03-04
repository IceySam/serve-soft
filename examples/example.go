package examples

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/IceySam/serve-soft/network"
)

var foods = map[int]string{
	0: "yam",
	1: "beans",
	2: "rice",
	3: "oil",
	4: "beef",
}
var resp network.Responses

func add(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		resp.RepondBadRequest(w, r, "Method not allowed")
	} else {
		f := &Food{}
		err := json.NewDecoder(r.Body).Decode(f)
		if err != nil {
			resp.RepondBadRequest(w, r, "Could not parse json data")
		} else if err := f.validate(); err != nil {
			resp.RepondBadRequest(w, r, err.Error())
		} else {
			foods[len(foods)] = f.Name
			resp.RespondCreated(w, r, f, f.Name)
		}
	}
}

func getOne(w http.ResponseWriter, r *http.Request) {
	last := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(last[2])

	if err != nil {
		resp.RepondBadRequest(w, r, err.Error())
	} else if id > -1 && id < len(foods) {
		data := foods[id]
		resp.RespondOk(w, r, data)
	} else {
		resp.RepondBadRequest(w, r, fmt.Sprintf("out of range: %d of %d", id, len(foods)))
	}
}

func getAll(w http.ResponseWriter, r *http.Request) {
	values := make([]string, 0, len(foods))
	for _, value := range foods {
		values = append(values, value)
	}

	data := make(map[string]any)
	data["foods"] = values
	resp.RespondOk(w, r, data)
}

func update(w http.ResponseWriter, r *http.Request) {
	last := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(last[3])
	if err != nil {
		resp.RepondBadRequest(w, r, err.Error())
	} else if r.Method != http.MethodPut {
		resp.RepondBadRequest(w, r, "Method not allowed")
	} else if id > -1 && id < len(foods) {
		f := &Food{}
		err := json.NewDecoder(r.Body).Decode(f)
		if err != nil {
			resp.RepondBadRequest(w, r, "Could not parse json data")
		} else if err := f.validate(); err != nil {
			resp.RepondBadRequest(w, r, err.Error())
		} else {
			foods[id] = f.Name
			resp.RespondUpdated(w, r)
		}
	} else {
		resp.RepondBadRequest(w, r, fmt.Sprintf("out of range: %d of %d", id, len(foods)))
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		resp.RepondBadRequest(w, r, "Method not allowed")
	} else {
		u := &User{}
		err := json.NewDecoder(r.Body).Decode(u)

		if err != nil {
			resp.RepondBadRequest(w, r, "Could not parse json data")
		} else if err := u.validate(); err != nil {
			resp.RepondBadRequest(w, r, err.Error())
		} else {
			_, token, err := network.GenerateClaim(u.Id, u.FirstName, u.OtherNames, "test", 1800)
			if err != nil {
				resp.RepondBadRequest(w, r, err.Error())
			} else {
				data := make(map[string]any)
				data["user"] = u
				data["token"] = token
				resp.RespondCreated(w, r, data, "Login successful")
			}
		}
	}
}

func Setup(h *network.NetHandler) {
	// initiate responses
	resp = network.Responses{}

	// auth middleware
	auth := h.Middlewares["auth"]

	h.Mux.HandleFunc("/auth/login", login) // should be in auth package
	h.Mux.HandleFunc("/foods/add", add)
	h.Mux.HandleFunc("/foods/update/", update)
	h.Mux.Handle("/foods", auth(http.HandlerFunc(getAll)))
	h.Mux.HandleFunc("/foods/", getOne)

}
