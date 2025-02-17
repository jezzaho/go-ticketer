package handler

import (
	"fmt"
	"net/http"
)

type Ticket struct{}

func (o *Ticket) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create a ticket")
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
