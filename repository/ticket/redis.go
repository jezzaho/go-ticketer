package ticket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jezzaho/go-ticketer/model"
	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	Client *redis.Client
}

func ticketIDKey(id uint64) string {
	return fmt.Sprintf("ticket:%d", id)
}

func (r *RedisRepo) Insert(ctx context.Context, ticket model.Ticket) error {
	data, err := json.Marshal(ticket)
	if err != nil {
		return fmt.Errorf("failed to encode ticket: %w", err)
	}

	key := ticketIDKey(ticket.TicketID)

	txn := r.Client.TxPipeline()

	res := txn.SetNX(ctx, key, string(data), 0)
	if err := res.Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to set: %w", err)
	}
	if err := txn.SAdd(ctx, "tickets", key).Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to add to tickets set: %w", err)
	}
	if _, err := txn.Exec(ctx); err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}

	return nil
}

var ErrNotExists = errors.New("order does not exist")

func (r *RedisRepo) FindByID(ctx context.Context, id uint64) (model.Ticket, error) {
	key := ticketIDKey(id)

	value, err := r.Client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return model.Ticket{}, ErrNotExists
	} else if err != nil {
		return model.Ticket{}, fmt.Errorf("got ticket: %w", err)
	}

	var ticket model.Ticket
	err = json.Unmarshal([]byte(value), &ticket)
	if err != nil {
		fmt.Errorf("failed to decode ticket json: %w", err)
	}

	return ticket, nil
}

func (r *RedisRepo) DeleteByID(ctx context.Context, id uint64) error {
	key := ticketIDKey(id)

	txn := r.Client.TxPipeline()

	err := txn.Del(ctx, key).Err()
	if errors.Is(err, redis.Nil) {
		txn.Discard()
		return ErrNotExists
	} else if err != nil {
		txn.Discard()
		return fmt.Errorf("get ticket: %w", err)
	}
	if err := txn.SRem(ctx, "tickets", key).Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to remove from tickets set: %w", err)
	}
	if _, err := txn.Exec(ctx); err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}

	return nil
}

func (r *RedisRepo) Update(ctx context.Context, ticket model.Ticket) error {
	data, err := json.Marshal(ticket)
	if err != nil {
		return fmt.Errorf("failed to encode order: %w", err)
	}

	key := ticketIDKey(ticket.TicketID)

	err = r.Client.SetXX(ctx, key, string(data), 0).Err()
	if errors.Is(err, redis.Nil) {
		return ErrNotExists
	}

	return nil
}

// Finding with pagination

type FindAllPage struct {
	Size   uint64
	Offset uint64
}

type FindResult struct {
	Tickets []model.Ticket
	Cursor  uint64
}

func (r *RedisRepo) FindAll(ctx context.Context, page FindAllPage) (FindResult, error) {
	res := r.Client.SScan(ctx, "tickets", page.Offset, "*", int64(page.Size))

	keys, cursor, err := res.Result()
	if err != nil {
		return FindResult{}, fmt.Errorf("failed to get tickets ids: %w", err)
	}
	if len(keys) == 0 {
		return FindResult{
			Tickets: []model.Ticket{},
		}, nil
	}

	xs, err := r.Client.MGet(ctx, keys...).Result()
	if err != nil {
		return FindResult{}, fmt.Errorf("failed to get tickets: %w", err)
	}
	tickets := make([]model.Ticket, len(xs))

	for i, x := range xs {
		x := x.(string)
		var ticket model.Ticket

		err := json.Unmarshal([]byte(x), &ticket)
		if err != nil {
			return FindResult{}, fmt.Errorf("failed to decode ticket json: %w", err)
		}
		tickets[i] = ticket
	}
	return FindResult{
		Tickets: tickets,
		Cursor:  cursor,
	}, nil
}
