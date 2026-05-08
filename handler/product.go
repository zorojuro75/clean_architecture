package handler

import (
    "net/http"
    "strconv"

    "productmanager/entity"
    "github.com/gin-gonic/gin"
)

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

type ProductHandler struct {
    uc entity.ProductUsecase
}

func NewProductHandler(uc entity.ProductUsecase) *ProductHandler {
    return &ProductHandler{uc: uc}
}



func parseID(c *gin.Context) (uint, bool) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return 0, false
    }
    return uint(id), true
}

func (h *ProductHandler) ListProducts(c *gin.Context) {
    page,  _ := strconv.Atoi(c.DefaultQuery("page",  "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

    products, total := h.uc.ListProducts(page, limit)
    respondPaginated(c, products, total, page, limit)
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
    id, ok := parseID(c)
    if !ok { return }

    product, err := h.uc.GetProduct(id)
    if err != nil { mapErr(c, err); return }
    responseOk(c, product)
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
    var req createProductReq
    err := c.ShouldBindJSON(&req);
    if err != nil {
        mapErr(c, err)
        return
    }

    product, err := h.uc.CreateProduct(req.Name, req.Price, req.Stock)
    if err != nil { mapErr(c, err); return }
    respondCreated(c, product)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
    id, ok := parseID(c)
    if !ok { return }

    var req updateProductReq
    err := c.ShouldBindJSON(&req)
    if err != nil {
        mapErr(c, err)
        return
    }

    product, err := h.uc.UpdateProduct(id, req.Name, req.Price, req.Stock)
    if err != nil {
        mapErr(c, err);
        return
    }
    respondCreated(c, product)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
    id, ok := parseID(c)
    if !ok { return }

    if err := h.uc.DeleteProduct(id); err != nil {
        mapErr(c, err)
        return
    }
    respondMessage(c, "product deleted")
}