package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
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

func (t *Ticket) List(w http.ResponseWriter, r *http.Request) {
	cursorStr := r.URL.Query().Get("cursor")
	if cursorStr == "" {
		cursorStr = "0"
	}

	const decimal = 10
	const bitSize = 64
	cursor, err := strconv.ParseUint(cursorStr, decimal, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	const size = 50
	res, err := t.Repo.FindAll(r.Context(), ticket.FindAllPage{
		Offset: cursor,
		Size:   size,
	})
	if err != nil {
		fmt.Println("failed to find all:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var response struct {
		Items []model.Ticket `json:"items"`
		Next  uint64         `json:"next,omitempty"`
	}
	response.Items = res.Tickets
	response.Next = res.Cursor

	data, err := json.Marshal(response)
	if err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(data)

}

func (t *Ticket) GetByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64

	ticketID, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	o, err := t.Repo.FindByID(r.Context(), ticketID)
	if errors.Is(err, ticket.ErrNotExists) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("failed to find by id: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(o); err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (t *Ticket) UpdateByID(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		fmt.Println("error decoding body:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64

	ticketId, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	theTicket, err := t.Repo.FindByID(r.Context(), ticketId)
	if errors.Is(err, ticket.ErrNotExists) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("failed to find by id: %w", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	now := time.Now().UTC()

	theTicket.Status = body.Status
	theTicket.LastUpdate = &now

	err = t.Repo.Update(r.Context(), theTicket)
	if err != nil {
		fmt.Println("failed to insert:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(theTicket); err != nil {
		fmt.Println("failed to marshal", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (t *Ticket) DeleteByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64

	ticketId, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = t.Repo.DeleteByID(r.Context(), ticketId)
	if errors.Is(err, ticket.ErrNotExists) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("failed to find by id:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
