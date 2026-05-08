package main

import (
	"productmanager/handler"
	"productmanager/repository"
	"productmanager/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	repo := repository.NewInMemoryRepo()
	productUC := usecase.NewProductUsecase(repo)
	productHandler := handler.NewProductHandler(productUC)

	r := gin.Default()
	healthHandler := handler.NewHealthHandler()
	r.GET("/health", healthHandler.Check)

	api := r.Group("/api/v1")
	{
		p := api.Group("/products")
		{
			p.GET("", productHandler.ListProducts)
			p.GET("/:id", productHandler.GetProduct)
			p.POST("", productHandler.CreateProduct)
			p.PUT("/:id", productHandler.UpdateProduct)
			p.DELETE("/:id", productHandler.DeleteProduct)
		}
	}

	r.Run(":8080")
}
