package shared

import (
	"strconv"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

const requestIDHeader = "X-Request-Id"
const requestIDContextKey = "request_id"
var requestIDCounter uint64

func requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(requestIDHeader)
		if requestID == "" {
			requestID = newRequestID()
		}
		c.Set(requestIDContextKey, requestID)
		c.Writer.Header().Set(requestIDHeader, requestID)
		c.Next()
	}
}

func GetRequestID(c *gin.Context) string {
	if value, ok := c.Get(requestIDContextKey); ok {
		if requestID, ok := value.(string); ok {
			return requestID
		}
	}
	return ""
}

func newRequestID() string {
	counter := atomic.AddUint64(&requestIDCounter, 1)
	timestamp := time.Now().UnixNano()
	return strconv.FormatInt(timestamp, 36) + "-" + strconv.FormatUint(counter, 36)
}
