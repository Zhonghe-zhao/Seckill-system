package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Zhonghe-zhao/seckill-system/internal/config"
	"github.com/Zhonghe-zhao/seckill-system/internal/handler"
	"github.com/Zhonghe-zhao/seckill-system/internal/model"
	"github.com/Zhonghe-zhao/seckill-system/internal/repository/db"
	"github.com/Zhonghe-zhao/seckill-system/internal/repository/rds"
	"github.com/Zhonghe-zhao/seckill-system/internal/router"
	"github.com/Zhonghe-zhao/seckill-system/internal/service"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.LoadConfig("../..")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Println("配置加载成功")

	// 初始化数据库
	pg, err := gorm.Open(postgres.Open(cfg.DBSource), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	RunMigrations(pg)
	// 初始化 Redis
	redisClient := InitRedisClient(&cfg)
	//初始化各层组件
	// 假设 service.ProductService 是一个结构体类型，需要先实例化它
	RedisRepository := rds.NewRedisRepository(redisClient)
	DRepository := db.NewDBRepository(pg)
	productService := service.NewProductService(RedisRepository, DRepository)
	productHandler := handler.NewProductHandler(productService)

	// 初始化路由器（你可以将 db 和 rdb 注入进去，如果后续需要使用）

	router := router.SetupRouter(productHandler)
	fmt.Printf("服务器启动成功，监听地址: %s\n", cfg.HTTPServerAddress)
	err = http.ListenAndServe(cfg.HTTPServerAddress, router)
	if err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}

}

func RunMigrations(db *gorm.DB) {
	log.Println("正在运行数据库迁移...")
	err := db.AutoMigrate(&model.Product{}, &model.Order{}, &model.User{})
	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}
	log.Println("数据库迁移完成!")
}

func InitRedisClient(cfg *config.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddress, // Redis 地址，如 "localhost:6379"
		DB:   0,                // 使用第 0 个逻辑数据库（Redis 有多个逻辑DB，默认是0）
	})
	pingCtx, cancelPing := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelPing()
	if _, err := rdb.Ping(pingCtx).Result(); err != nil {
		log.Fatalf("无法连接到 Redis: %v", err)
	}
	log.Println("成功连接到 Redis!")
	return rdb
}
