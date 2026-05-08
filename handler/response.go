package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"errors"
	"productmanager/entity"
)

type successResponse struct {
	Success bool 	`json:"success"`
	Data 	any 	`json:"data,omitempty"`
	Message any		`json:"message,omitempty"`
}

type errorResponse struct {
	Success bool	`json:"success"`
	Error 	string	`json:"error"`
}

type paginatedResponse struct {
	Success bool 	`json:"success"`
	Data 	any		`json:"data"`
	Total	int 	`json:"total"`
	Page 	int		`json:"page"`
	Limit 	int 	`json:"limit"`
}

func responseOk(c *gin.Context, data any){
	c.JSON(http.StatusOK, successResponse{Success: true, Data: data})
}

func respondCreated(c *gin.Context, data any) {
    c.JSON(http.StatusCreated, successResponse{Success: true, Data: data})
}

func respondMessage(c *gin.Context, msg string) {
    c.JSON(http.StatusOK, successResponse{Success: true, Message: msg})
}

func respondPaginated(c *gin.Context, data any, total, page, limit int) {
    c.JSON(http.StatusOK, paginatedResponse{
        Success: true, Data: data,
        Total: total, Page: page, Limit: limit,
    })
}

func respondError(c *gin.Context, status int, msg string) {
    c.JSON(status, errorResponse{Success: false, Error: msg})
}

func mapErr(c *gin.Context, err error) {
    switch {
    case errors.Is(err, entity.ErrNotFound):
        respondError(c, http.StatusNotFound, err.Error())
    case errors.Is(err, entity.ErrInvalidInput):
        respondError(c, http.StatusBadRequest, err.Error())
    case errors.Is(err, entity.ErrDuplicate):
        respondError(c, http.StatusConflict, err.Error())
    case errors.Is(err, entity.ErrOutOfStock):
        respondError(c, http.StatusConflict, err.Error())
    default:
        respondError(c, http.StatusInternalServerError, "internal server error")
    }
}