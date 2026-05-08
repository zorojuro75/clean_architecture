package handler

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
)

var startTime = time.Now()

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
    return &HealthHandler{}
}

func (h *HealthHandler) Check(c *gin.Context) {
    uptime := time.Since(startTime).Round(time.Second)

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "status":  "ok",
        "version": "1.0.0",
        "uptime":  uptime.String(),
    })
}