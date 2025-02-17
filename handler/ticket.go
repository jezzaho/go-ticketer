package handler

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/jezzaho/go-ticketer/model"
	"github.com/jezzaho/go-ticketer/repository/ticket"
)

type Ticket struct {
	Repo *ticket.RedisRepo
}

func (t *Ticket) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		// CreatorID   uuid.UUID `json:"creator_id"`
		Labels      []string `json:"labels"`
		Status      string   `json:"status"`
		Category    string   `json:"category"`
		Description string   `json:"description"`
		// Add DueTo field - for easier now it is a deadhead
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()

	ticket := model.Ticket{
		TicketID:    rand.Uint64(),
		CreatorID:   rand.Uint64(),
		Labels:      body.Labels,
		Status:      body.Status,
		Category:    body.Category,
		CreatedAt:   &now,
		LastUpdate:  &now,
		DueTo:       &now,
		Description: body.Description,
	}

	err := t.Repo.Insert(r.Context(), ticket)
	if err != nil {
		fmt.Println("failed to insert:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	res, err := json.Marshal(ticket)
	if err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

func (o *Ticket) List(w http.ResponseWriter, r *http.Request) {
	fmt.Println("List all tickets")
}

func (o *Ticket) GetByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get a ticket by ID")
}

func (o *Ticket) UpdateByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update a ticket by ID")
}

func (o *Ticket) DeleteByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete a ticket by ID")
}
