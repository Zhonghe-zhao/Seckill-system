package main

import (
	"log"

	"github.com/Zhonghe-zhao/seckill-system/internal/config"
)

func main() {
	cfg := config.LoadConfig()
	r := router.SetupRouter(cfg)

	err := r.Run(cfg.ServerAddress)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
