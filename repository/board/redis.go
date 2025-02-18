package board

import (
	"context"
	"fmt"

	"github.com/jezzaho/go-ticketer/model"
	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	Client *redis.Client
}

func boardIDKey(id uint64) string {
	return fmt.Sprintf("ticket:%d", id)
}

func (r *RedisRepo) Insert(ctx context.Context, board model.Board) error {
	return fmt.Errorf("empty")
}
func (r *RedisRepo) FindByID(ctx context.Context, id uint64) (model.Board, error) {
	return model.Board{}, fmt.Errorf("empty")
}
func (r *RedisRepo) DeleteByID(ctx context.Context, id uint64) error {
	return fmt.Errorf("empty")
}
func (r *RedisRepo) Update(ctx context.Context, board model.Board) error {
	return fmt.Errorf("empty")
}
func (r *RedisRepo) FindAll(ctx context.Context) ([]model.Board, error) {
	return nil, fmt.Errorf("empty")
}
