package handler

import (
	"net/http"

	"github.com/jezzaho/go-ticketer/repository/board"
)

type Board struct {
	Repo *board.RedisRepo
}

func (b *Board) Create(w http.ResponseWriter, r *http.Request) {

}

func (b *Board) List(w http.ResponseWriter, r *http.Request) {

}

func (b *Board) Update(w http.ResponseWriter, r *http.Request) {

}
func (b *Board) Delete(w http.ResponseWriter, r *http.Request) {

}

func (b *Board) GetByID(w http.ResponseWriter, r *http.Request) {

}
