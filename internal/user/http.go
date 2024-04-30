package user

import (
	"encoding/json"
	"log"
	"net/http"
)

type Handler struct {
	Svc service
}

func (h Handler) AddUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var u User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request body"))
		return
	}
	// Call the AddUser function
	message, err := h.Svc.AddUser(u)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to add user"))
		return
	}

	// Return a success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
}
