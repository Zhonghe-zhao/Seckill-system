package service

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Zhonghe-zhao/seckill-system/internal/model"
	"github.com/Zhonghe-zhao/seckill-system/internal/repository"
)

type ProductService struct {
	redisStockRepo       repository.StockRepository       // Redis 库存操作
	redisProductInfoRepo repository.ProductInfoRepository // Redis 商品信息缓存操作
	dbRepo               repository.DBRepository          // GORM 数据库商品操作
}

func NewProductService(
	rsr repository.StockRepository,
	rpir repository.ProductInfoRepository,
	dbr repository.DBRepository,
) *ProductService {
	return &ProductService{
		redisStockRepo:       rsr,
		redisProductInfoRepo: rpir,
		dbRepo:               dbr,
	}
}

func (s *ProductService) InitializeProduct(ctx context.Context, name, description string, price float64, initialStock int64, startTime, endTime time.Time) (*model.Product, error) {
	dbProduct := &model.Product{
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       initialStock,
		StartTime:   startTime,
		EndTime:     endTime,
	}
	if err := s.dbRepo.CreateProduct(ctx, dbProduct); err != nil {
		return nil, fmt.Errorf("创建商品到数据库失败: %w", err)
	}
	log.Printf("商品 '%s' (DB ID: %d) 成功保存到数据库", dbProduct.Name, dbProduct.ID)

	redisProductIDStr := strconv.FormatUint(uint64(dbProduct.ID), 10)

	// 缓存商品信息到 Redis 中 这里的 model.Product 是用于缓存的结构，其 Stock 字段应反映初始库存或当前可售状态
	cacheProduct := &model.Product{ // 创建一个新的实例用于缓存，避免直接修改 dbProduct
		Name:        dbProduct.Name,
		Description: dbProduct.Description, // 可能不需要缓存完整描述
		Price:       dbProduct.Price,
		Stock:       initialStock, // 当前可售库存（刚初始化时与初始库存一致）
		StartTime:   dbProduct.StartTime,
		EndTime:     dbProduct.EndTime,
	}
	if err := s.redisProductInfoRepo.SetProductInfo(ctx, cacheProduct); err != nil {
		log.Printf("警告: 商品 (DB ID: %d) 数据库创建成功，但缓存商品信息到 Redis 失败: %v. Key: product_info:%s", dbProduct.ID, err, redisProductIDStr)
	}
	log.Printf("商品 '%s' (DB ID: %d, Redis Key ID: %s) 初始化完成，初始库存: %d", dbProduct.Name, dbProduct.ID, redisProductIDStr, initialStock)
	return dbProduct, nil // 返回从数据库创建并包含自增 ID 的对象
}
