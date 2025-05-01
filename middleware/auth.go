package middleware

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func RequireAuth(role string) gin.HandlerFunc {
    return func(c *gin.Context) {
        user, exists := c.Get("user")
        if !exists || user.(string) != role {
            c.Redirect(http.StatusSeeOther, "/login")
            c.Abort()
            return
        }
        c.Next()
    }
}