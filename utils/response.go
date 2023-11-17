package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Status  int
	Message []string
	Error   []string
}

func (r *Response) Success(c *gin.Context, msg interface{}) {
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": msg,
	})
}
func (r *Response) BadRequest(c *gin.Context, msg interface{}) {
	c.Header("Content-Type", "application/json")
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"status":  http.StatusBadRequest,
		"message": msg,
	})
}
