package main

import (
    "productmanager/handler"
    "productmanager/repository"
    "productmanager/usecase"

    "github.com/gin-gonic/gin"
)

func main() {
    // 1. Infrastructure
    repo := repository.NewInMemoryRepo()

    // 2. Business logic
    productUC := usecase.NewProductUsecase(repo)

    // 3. Delivery — inject usecase into handler
    productHandler := handler.NewProductHandler(productUC)

    // 4. Router
    r := gin.Default()

    api := r.Group("/api/v1")
    {
        p := api.Group("/products")
        {
            p.GET("",     productHandler.ListProducts)
            p.GET("/:id", productHandler.GetProduct)
            p.POST("",    productHandler.CreateProduct)
            p.PUT("/:id", productHandler.UpdateProduct)
            p.DELETE("/:id", productHandler.DeleteProduct)
        }
    }

    r.Run(":8080")
}