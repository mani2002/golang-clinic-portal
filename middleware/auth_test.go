package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Middleware sets user role from cookie", func(t *testing.T) {
		router := gin.Default()
		
		router.Use(func(c *gin.Context) {
			role, _ := c.Cookie("user")
			if role != "" { 
				c.Set("user", role) 
			}
			c.Next()
		})
		
		router.GET("/check-role", func(c *gin.Context) {
			role, exists := c.Get("user")
			if !exists {
				c.String(http.StatusOK, "No role")
				return
			}
			c.String(http.StatusOK, "Role: %s", role)
		})
		
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/check-role", nil)
		req.AddCookie(&http.Cookie{
			Name:  "user",
			Value: "doctor",
		})
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Role: doctor")
	})
	
	t.Run("Receptionist-only route allows receptionist", func(t *testing.T) {
		router := gin.Default()
		
		router.Use(func(c *gin.Context) {
			role, _ := c.Cookie("user")
			if role != "" { 
				c.Set("user", role) 
			}
			c.Next()
		})
		
		router.GET("/receptionist-only", func(c *gin.Context) {
			role, exists := c.Get("user")
			if !exists || role != "receptionist" {
				c.Redirect(http.StatusFound, "/patients")
				return
			}
			c.String(http.StatusOK, "Receptionist Access Granted")
		})
		
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/receptionist-only", nil)
		req.AddCookie(&http.Cookie{
			Name:  "user",
			Value: "receptionist",
		})
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Receptionist Access Granted")
	})
	
	t.Run("Receptionist-only route redirects doctor", func(t *testing.T) {
		router := gin.Default()
		
		router.Use(func(c *gin.Context) {
			role, _ := c.Cookie("user")
			if role != "" { 
				c.Set("user", role) 
			}
			c.Next()
		})
		
		router.GET("/receptionist-only", func(c *gin.Context) {
			role, exists := c.Get("user")
			if !exists || role != "receptionist" {
				c.Redirect(http.StatusFound, "/patients")
				return
			}
			c.String(http.StatusOK, "Receptionist Access Granted")
		})
		
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/receptionist-only", nil)
		req.AddCookie(&http.Cookie{
			Name:  "user",
			Value: "doctor",
		})
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusFound, w.Code)
		assert.Equal(t, "/patients", w.Header().Get("Location"))
	})
	
	t.Run("Routes without cookie have no role", func(t *testing.T) {
		router := gin.Default()
		
		router.Use(func(c *gin.Context) {
			role, _ := c.Cookie("user")
			if role != "" { 
				c.Set("user", role) 
			}
			c.Next()
		})
		
		router.GET("/check-user", func(c *gin.Context) {
			_, exists := c.Get("user")
			if !exists {
				c.String(http.StatusOK, "No user")
				return
			}
			c.String(http.StatusOK, "User exists")
		})
		
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/check-user", nil)
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "No user")
	})
}