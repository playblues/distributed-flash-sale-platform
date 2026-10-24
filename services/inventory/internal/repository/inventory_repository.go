package repository

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type InventoryRepository struct {
	rdb           *redis.Client
	reserveScript *redis.Script
}

func NewInventoryRepository(rdb *redis.Client) *InventoryRepository {
	luaScript := redis.NewScript(`
		if redis.call("exists", KEYS[2]) == 1 then
			return -1
		end
		local stock = tonumber(redis.call("get", KEYS[1]))
	
		if stock and stock > 0 then
			redis.call('decr', KEYS[1])
			redis.call('setex', KEYS[2], 600, 'reserved')
			
			return 1
		else
			return 0
		end
	`)
	return &InventoryRepository{
		rdb:           rdb,
		reserveScript: luaScript,
	}
}
func (r *InventoryRepository) InitializeStock(ctx context.Context,
	eventID string, amount int) error {
	key := fmt.Sprintf("inventory:event:%s", eventID)
	return r.rdb.Set(ctx, key, amount, 0).Err()
}

func (r *InventoryRepository) GetStock(ctx context.Context, eventID string) (int, error) {
	key := fmt.Sprintf("inventory:event:%s", eventID)
	val, err := r.rdb.Get(ctx, key).Int()
	if err == redis.Nil {
		return 0, nil
	}
	return val, err
}
func (r *InventoryRepository) ReserveStock(ctx context.Context,
	eventID, userID string) (int64, error) {

	stockKey := fmt.Sprintf("inventory:event:%s", eventID)
	reserveKey := fmt.Sprintf("reserve:user:%s:event:%s", userID, eventID)

	result, err := r.reserveScript.Run(ctx, r.rdb, []string{stockKey, reserveKey}).Result()

	if err != nil {
		return 0, err
	}
	return result.(int64), nil
}
func (r *InventoryRepository) RollBack(ctx context.Context, eventID, userId string) error {
	stockKey := fmt.Sprintf("inventory:event:%s", eventID)
	reserveKey := fmt.Sprintf("reserve:user:%s:event:%s", userId, eventID)

	pipe := r.rdb.TxPipeline()
	pipe.Incr(ctx, stockKey)
	pipe.Del(ctx, reserveKey)

	_, err := pipe.Exec(ctx)
	return err
}
