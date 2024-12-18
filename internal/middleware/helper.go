package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func CreateErrResp(format string, a ...any) interface{} {
	return gin.H{"error": fmt.Sprintf(format, a...)}
}
