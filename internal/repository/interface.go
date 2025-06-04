package repository

import (
	"context"
	"github.com/Zhonghe-zhao/seckill-system/internal/model"
	"github.com/Zhonghe-zhao/seckill-system/internal/repository/db"
	"github.com/Zhonghe-zhao/seckill-system/internal/repository/rds"
)

// Redis 商品缓存
type ProductCache interface {
	GetProductInfo(ctx context.Context, productID string) (*model.Product, error)
	SetProductInfo(ctx context.Context, product *model.Product) error
	GetStock(ctx context.Context, productID string) (int64, error)
	SetStock(ctx context.Context, productID string, stock int64) error
	DecrementStock(ctx context.Context, productID string) (int64, error)
}

// Redis 订单缓存逻辑（是否下过单、入队）
type OrderCache interface {
	AddUserToPurchasedSet(ctx context.Context, productID, userID string) (bool, error)
	IsUserInPurchasedSet(ctx context.Context, productID, userID string) (bool, error)
	PushOrderToQueue(ctx context.Context, queueName string, orderReq *model.OrderRequest) error
}

// DB 商品持久化
type ProductStorage interface {
	GetProductByID(ctx context.Context, productID uint) (*model.Product, error)
	DecrementStockInDB(ctx context.Context, productID uint) error
	CreateProduct(ctx context.Context, product *model.Product) error
}

// DB 订单持久化

type OrderStorage interface {
	CreateOrder(ctx context.Context, order *model.Order) error
	HasUserOrdered(ctx context.Context, userID uint, productID uint) (bool, error)
}

var _ ProductStorage = (*db.DBRepository)(nil)
var _ ProductCache = (*rds.RedisRepository)(nil)
