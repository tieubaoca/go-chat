package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tieubaoca/go-chat-server/dto/response"
	"github.com/tieubaoca/go-chat-server/types"
	"github.com/tieubaoca/go-chat-server/utils"
	"golang.org/x/time/rate"
)

type LimiterMiddleware struct {
	limiterCache map[string]*rate.Limiter
	whitelist    []string
	maxRequest   int
	duration     time.Duration
}

func NewLimiterMiddleware(
	whitelist []string,
	maxRequest int,
	duration time.Duration,
) *LimiterMiddleware {
	return &LimiterMiddleware{
		limiterCache: make(map[string]*rate.Limiter),
		whitelist:    whitelist,
		maxRequest:   maxRequest,
		duration:     duration,
	}
}

func (m *LimiterMiddleware) IPRateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if utils.ContainsString(m.whitelist, ip) {
			c.Next()
			return
		}
		limiter, exists := m.limiterCache[ip]
		if !exists {
			limiter = rate.NewLimiter(rate.Every(m.duration/time.Duration(m.maxRequest)), m.maxRequest)
			m.limiterCache[ip] = limiter
		}

		if !limiter.Allow() {
			c.JSON(
				http.StatusTooManyRequests,
				response.ResponseData{
					Status:  types.StatusError,
					Message: "Too many requests",
					Data:    "",
				},
			)
			c.Abort()
			return
		}

		c.Next()
	}
}
