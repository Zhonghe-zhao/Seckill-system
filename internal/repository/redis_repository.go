package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Zhonghe-zhao/seckill-system/internal/model"
	"github.com/redis/go-redis/v9"
)

type redisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) *redisRepository {
	return &redisRepository{client: client}
}

func (r *redisRepository) GetStockKey(productID string) string {
	return fmt.Sprintf("stock:%s", productID)
}

func (r *redisRepository) GetStock(ctx context.Context, productID string) (int64, error) {
	val, err := r.client.Get(ctx, r.GetStockKey(productID)).Int64()
	if err == redis.Nil {
		return 0, errors.New("product stock not found") // 或者返回特定错误类型
	} else if err != nil {
		return 0, err
	}
	return val, nil
}

func (r *redisRepository) SetStock(ctx context.Context, productID string, stock int64) error {
	return r.client.Set(ctx, r.GetStockKey(productID), stock, 0).Err() // 0 表示不过期
}

func (r *redisRepository) DecrementStock(ctx context.Context, productID string) (int64, error) {
	// 注意：这种方式在 DECR 后库存 < 0 时需要补偿逻辑，Lua 脚本更优
	return r.client.Decr(ctx, r.GetStockKey(productID)).Result()
}

// Lua 脚本定义
// func (r *redisRepository) DecrementStockLua(ctx context.Context, productID string) (int64, error) {
// 	keys := []string{r.GetStockKey(productID)}
// 	result, err := decrementStockScript.Run(ctx, r.client, keys).Result()
// 	if err != nil {
// 		return -3, fmt.Errorf("lua script error: %w", err)
// 	}
// 	return result.(int64), nil
// }

// --- OrderRepository
func (r *redisRepository) GetPurchasedSetKey(productID string) string {
	return fmt.Sprintf("purchased_users:%s", productID)
}

func (r *redisRepository) AddUserToPurchasedSet(ctx context.Context, productID string, userID string) (bool, error) {
	// SADD 返回成功添加的元素数量，1 表示新添加，0 表示已存在
	added, err := r.client.SAdd(ctx, r.GetPurchasedSetKey(productID), userID).Result()
	if err != nil {
		return false, err
	}
	return added == 1, nil
}

func (r *redisRepository) IsUserInPurchasedSet(ctx context.Context, productID string, userID string) (bool, error) {
	return r.client.SIsMember(ctx, r.GetPurchasedSetKey(productID), userID).Result()
}

func (r *redisRepository) GetOrderQueueKey() string {
	return "order_queue" // 简化，所有商品共用一个队列
}

func (r *redisRepository) PushOrderToQueue(ctx context.Context, queueName string, orderReq *model.OrderRequest) error {
	// 实际项目中，orderReq 需要序列化，例如 JSON
	// 这里简化为直接推入一个代表性的字符串，或者你可以传递序列化后的字节流
	orderData := fmt.Sprintf("%s:%s:%d", orderReq.UserID, orderReq.ProductID, orderReq.Timestamp)
	return r.client.LPush(ctx, queueName, orderData).Err()
}

// --- ProductInfoRepository

func (r *redisRepository) GetProductInfoKey(productID string) string {
	return fmt.Sprintf("product_info:%s", productID)
}

// GetProductInfo 示例 (实际应存储和检索更多字段)
func (r *redisRepository) GetProductInfo(ctx context.Context, productID string) (*model.Product, error) {
	// 假设将整个 Product struct 以 JSON 字符串形式存储在 Redis String 中
	// 或者使用 Redis Hash 存储各个字段
	data, err := r.client.Get(ctx, r.GetProductInfoKey(productID)).Bytes()
	if err == redis.Nil {
		return nil, errors.New("product info not found in cache")
	} else if err != nil {
		return nil, err
	}
	// 这里需要反序列化 data 到 model.Product (例如使用 json.Unmarshal)
	// 此处简化，实际项目中你需要实现序列化/反序列化
	var product model.Product
	err = json.Unmarshal(data, &product)
	if err != nil {
		return nil, err
	}
	log.Printf("INFO: Product loaded from Redis: ID=%s", product.ID)

	return &product, nil
}

func (r *redisRepository) SetProductInfo(ctx context.Context, product *model.Product) error {
	data, err := json.Marshal(product)
	if err != nil {
		return err
	}

	log.Printf("INFO: SetProductInfo caching product: ID=%s", product.ID)

	return r.client.Set(ctx, r.GetProductInfoKey(strconv.FormatUint(uint64(product.ID), 10)), data, 24*time.Hour).Err()

}
