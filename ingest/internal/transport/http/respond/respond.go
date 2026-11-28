package respond

import "github.com/gin-gonic/gin"

type ErrorResponse struct {
	Error string `json:"error"`
}

func Error(c *gin.Context, status int, msg string) {
	c.AbortWithStatusJSON(status, ErrorResponse{Error: msg})
}
