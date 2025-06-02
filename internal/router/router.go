package router

import (
	"fmt"
	"net/http"

	"github.com/Zhonghe-zhao/seckill-system/internal/handler"
)

func SetupRouter(ph *handler.ProductHandler /*sh *handler.SaleHandler*/) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"status":"OK"}`)
	})

	//手动分组路由 因为net/http不支持分组路由
	mux.HandleFunc("POST /api/v1/admin/product/initialize", ph.HandleInitializeProduct)

	//mux.HandleFunc("GET /api/v1/product", ph.HandleGetProductDetails)
	//mux.HandleFunc("POST /api/v1/sale/attempt", sh.HandleSaleAttempt)
	// mux.HandleFunc("GET /api/v1/product", ph.HandleGetProductDetails)
	// mux.HandleFunc("POST /api/v1/sale/attempt", sh.HandleSaleAttempt)
	return mux
}
