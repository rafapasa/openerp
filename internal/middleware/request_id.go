package middleware

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid" // Import the uuid package
)

func RequestIDMiddleware() gin.HandlerFunc {
	return requestid.New(
		requestid.WithGenerator(func() string {
			v, err := uuid.NewUUID() // Use uuid.NewUUID() for V1 UUID
			if err != nil {
				return ""
			}
			return v.String()
		}),
		requestid.WithCustomHeaderStrKey("X-Request-ID"),
	)
}
