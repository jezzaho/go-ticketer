package model

import (
	"time"
)

type Ticket struct {
	TicketID    uint64     `json:"ticket_id"`
	CreatorID   uint64     `json:"creator_id`
	Labels      []string   `json:"labels"`
	Status      string     `json:"status"`
	Category    string     `json:"category"`
	CreatedAt   *time.Time `json:"created_at"`
	LastUpdate  *time.Time `json:"last_update"`
	DueTo       *time.Time `json:"due_to"`
	Description string     `json:"description"`
}
