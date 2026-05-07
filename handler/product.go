package handler

import (
    "errors"
    "net/http"
    "strconv"

    "productmanager/entity"
    "github.com/gin-gonic/gin"
)

// Request structs — Gin binding layer only
type createProductReq struct {
    Name  string  `json:"name"  binding:"required"`
    Price float64 `json:"price" binding:"required,gt=0"`
    Stock int     `json:"stock" binding:"min=0"`
}

type updateProductReq struct {
    Name  string  `json:"name"  binding:"required"`
    Price float64 `json:"price" binding:"required,gt=0"`
    Stock int     `json:"stock" binding:"min=0"`
}

// Handler struct — usecase injected
type ProductHandler struct {
    uc entity.ProductUsecase
}

func NewProductHandler(uc entity.ProductUsecase) *ProductHandler {
    return &ProductHandler{uc: uc}
}

// mapErr — converts domain errors to HTTP status codes
func mapErr(c *gin.Context, err error) {
    switch {
    case errors.Is(err, entity.ErrNotFound):
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
    case errors.Is(err, entity.ErrInvalidInput):
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    case errors.Is(err, entity.ErrDuplicate):
        c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
    case errors.Is(err, entity.ErrOutOfStock):
        c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
    default:
        c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
    }
}

// parseID — reusable URL param parser
func parseID(c *gin.Context) (uint, bool) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return 0, false
    }
    return uint(id), true
}

// GET /products
func (h *ProductHandler) ListProducts(c *gin.Context) {
    products := h.uc.ListProducts()
    c.JSON(http.StatusOK, gin.H{
        "data":  products,
        "total": len(products),
    })
}

// GET /products/:id
func (h *ProductHandler) GetProduct(c *gin.Context) {
    id, ok := parseID(c)
    if !ok { return }

    product, err := h.uc.GetProduct(id)
    if err != nil { mapErr(c, err); return }

    c.JSON(http.StatusOK, product)
}

// POST /products
func (h *ProductHandler) CreateProduct(c *gin.Context) {
    var req createProductReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    product, err := h.uc.CreateProduct(req.Name, req.Price, req.Stock)
    if err != nil { mapErr(c, err); return }

    c.JSON(http.StatusCreated, product)
}

// PUT /products/:id
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
    id, ok := parseID(c)
    if !ok { return }

    var req updateProductReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    product, err := h.uc.UpdateProduct(id, req.Name, req.Price, req.Stock)
    if err != nil { mapErr(c, err); return }

    c.JSON(http.StatusOK, product)
}

// DELETE /products/:id
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
    id, ok := parseID(c)
    if !ok { return }

    if err := h.uc.DeleteProduct(id); err != nil {
        mapErr(c, err)
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "product deleted"})
}