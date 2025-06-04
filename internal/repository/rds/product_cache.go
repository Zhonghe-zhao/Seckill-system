package rds

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{client: client}
}

func (r *RedisRepository) GetStockKey(productID string) string {
	return fmt.Sprintf("stock:%s", productID)
}

func (r *RedisRepository) GetStock(ctx context.Context, productID string) (int64, error) {
	val, err := r.client.Get(ctx, r.GetStockKey(productID)).Int64()
	if err == redis.Nil {
		return 0, errors.New("product stock not found") // 或者返回特定错误类型
	} else if err != nil {
		return 0, err
	}
	return val, nil
}

func (r *RedisRepository) SetStock(ctx context.Context, productID string, stock int64) error {
	return r.client.Set(ctx, r.GetStockKey(productID), stock, 0).Err() // 0 表示不过期
}

func (r *RedisRepository) DecrementStock(ctx context.Context, productID string) (int64, error) {
	// 注意：这种方式在 DECR 后库存 < 0 时需要补偿逻辑，Lua 脚本更优
	return r.client.Decr(ctx, r.GetStockKey(productID)).Result()
}

/*
// Lua 脚本定义
// func (r *RedisRepository) DecrementStockLua(ctx context.Context, productID string) (int64, error) {
// 	keys := []string{r.GetStockKey(productID)}
// 	result, err := decrementStockScript.Run(ctx, r.client, keys).Result()
// 	if err != nil {
// 		return -3, fmt.Errorf("lua script error: %w", err)
// 	}
// 	return result.(int64), nil
// }
*/
