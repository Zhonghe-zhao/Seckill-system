package repository

import (
	"context"

	"github.com/Zhonghe-zhao/seckill-system/internal/model"
	"gorm.io/gorm"
)

type dbRepository struct {
	db *gorm.DB
}

func NewDBRepository(db *gorm.DB) *dbRepository {
	return &dbRepository{db: db}
}

func (r *dbRepository) GetProductByID(ctx context.Context, productID uint) (*model.Product, error) {
	var product model.Product
	if err := r.db.WithContext(ctx).First(&product, productID).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *dbRepository) DecrementStock(ctx context.Context, productID uint) error {
	return r.db.WithContext(ctx).Model(&model.Product{}).
		Where("id = ? AND stock > 0", productID).
		Update("stock", gorm.Expr("stock - ?", 1)).Error
}

func (r *dbRepository) CreateOrder(ctx context.Context, order *model.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *dbRepository) HasUserOrdered(ctx context.Context, userID uint, productID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Order{}).
		Where("user_id = ? AND product_id = ?", userID, productID).
		Count(&count).Error
	return count > 0, err
}

func (r *dbRepository) CreateProduct(ctx context.Context, product *model.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}
