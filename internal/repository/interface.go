package repository

import (
	"context"

	"github.com/Zhonghe-zhao/seckill-system/internal/model"
)

// 负责操作 Redis 中库存
type StockRepository interface {
	GetStock(ctx context.Context, productID string) (int64, error)
	SetStock(ctx context.Context, productID string, stock int64) error
	DecrementStock(ctx context.Context, productID string) (int64, error)
	//DecrementStockLua(ctx context.Context, productID string) (int64, error)
}

// 负责处理订单相关 Redis 操作（用户是否下单、入队等）
type OrderRepository interface {
	AddUserToPurchasedSet(ctx context.Context, productID string, userID string) (bool, error)
	IsUserInPurchasedSet(ctx context.Context, productID string, userID string) (bool, error)
	PushOrderToQueue(ctx context.Context, queueName string, orderReq *model.OrderRequest) error
}

// Redis 缓存商品信息
type ProductInfoRepository interface {
	GetProductInfo(ctx context.Context, productID string) (*model.Product, error)
	SetProductInfo(ctx context.Context, product *model.Product) error
}

// 数据库操作接口
type DBRepository interface {
	GetProductByID(ctx context.Context, productID uint) (*model.Product, error)
	DecrementStock(ctx context.Context, productID uint) error
	//HasUserOrdered(ctx context.Context, userID uint, productID uint) (bool, error) //后期添加
	CreateProduct(ctx context.Context, product *model.Product) error

	CreateOrder(ctx context.Context, order *model.Order) error
}

var _ StockRepository = (*redisRepository)(nil)
var _ OrderRepository = (*redisRepository)(nil)
var _ ProductInfoRepository = (*redisRepository)(nil)
var _ DBRepository = (*dbRepository)(nil)
