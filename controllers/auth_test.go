package controllers

import (
	"github.com/mani2002/golang-clinic-portal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"fmt"
	"time"
)

func setupAuthTestDB(t *testing.T) *gorm.DB {
	dsn := "host=localhost user=postgres password=postgres dbname=clinicdb port=5432 sslmode=disable"
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	
	return db
}

func TestShowLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	t.Run("ShowLogin renders login page", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/login", nil)
		
		router := gin.Default()
		router.LoadHTMLGlob("../templates/*")
		router.GET("/login-test", ShowLogin)
		
		resp := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/login-test", nil)
		router.ServeHTTP(resp, req)
		
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), "Clinic Portal Login")
	})
}

func TestPostLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupAuthTestDB(t)
	
	timestamp := time.Now().UnixNano()
	testEmail := fmt.Sprintf("test_login_%d@example.com", timestamp)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	
	testUser := models.User{
		Email:    testEmail,
		Password: string(hashedPassword),
		Role:     "doctor",
	}
	db.Create(&testUser)
	
	t.Cleanup(func() {
		db.Unscoped().Where("email = ?", testEmail).Delete(&models.User{})
	})
	
	t.Run("PostLogin with valid credentials", func(t *testing.T) {
		router := gin.Default()
		router.LoadHTMLGlob("../templates/*")
		router.POST("/login", PostLogin(db))
		
		form := url.Values{}
		form.Add("email", testEmail)
		form.Add("password", "password123")
		form.Add("role", "doctor")
		
		resp := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/login", strings.NewReader(form.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		router.ServeHTTP(resp, req)
		
		assert.Equal(t, http.StatusSeeOther, resp.Code) // This should already be 303
		assert.Equal(t, "/patients", resp.Header().Get("Location"))
		
		cookies := resp.Result().Cookies()
		var userCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "user" {
				userCookie = cookie
				break
			}
		}
		assert.NotNil(t, userCookie)
		assert.Equal(t, "doctor", userCookie.Value)
	})
	
	t.Run("PostLogin with invalid credentials", func(t *testing.T) {
		router := gin.Default()
		router.LoadHTMLGlob("../templates/*")
		router.POST("/login", PostLogin(db))
		
		form := url.Values{}
		form.Add("email", testEmail)
		form.Add("password", "wrongpassword")
		form.Add("role", "doctor")
		
		resp := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/login", strings.NewReader(form.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		router.ServeHTTP(resp, req)
		
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), "Invalid credentials")
	})
	
	t.Run("PostLogin with non-existent user", func(t *testing.T) {
		router := gin.Default()
		router.LoadHTMLGlob("../templates/*")
		router.POST("/login", PostLogin(db))
		
		form := url.Values{}
		form.Add("email", "nonexistent@example.com")
		form.Add("password", "password123")
		form.Add("role", "receptionist")
		
		resp := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/login", strings.NewReader(form.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		router.ServeHTTP(resp, req)
		
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), "Invalid credentials")
	})
}