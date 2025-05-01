package controllers

import (
	"github.com/mani2002/golang-clinic-portal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"fmt"
	"time"
)

func setupPatientControllerTestDB(t *testing.T) *gorm.DB {
	dsn := "host=localhost user=postgres password=postgres dbname=clinicdb port=5432 sslmode=disable"
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	
	return db
}

func TestListPatients(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupPatientControllerTestDB(t)
	
	timestamp := time.Now().UnixNano()
	testPrefix := fmt.Sprintf("Test Patient %d", timestamp)
	
	testPatient1 := models.Patient{
		Name:   testPrefix + " - 1",
		Age:    30,
		Detail: "Test detail 1",
	}
	
	testPatient2 := models.Patient{
		Name:   testPrefix + " - 2",
		Age:    45,
		Detail: "Test detail 2",
	}
	
	db.Create(&testPatient1)
	db.Create(&testPatient2)
	
	t.Cleanup(func() {
		db.Unscoped().Where("name LIKE ?", testPrefix+"%").Delete(&models.Patient{})
	})
	
	t.Run("ListPatients displays patients list", func(t *testing.T) {
		router := gin.Default()
		router.LoadHTMLGlob("../templates/*")
		
		router.Use(func(c *gin.Context) {
			role, _ := c.Cookie("user")
			if role != "" { 
				c.Set("user", role) 
			}
			c.Next()
		})
		
		router.GET("/patients", ListPatients(db))
		
		resp := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/patients", nil)

		req.AddCookie(&http.Cookie{
			Name:  "user",
			Value: "doctor",
		})
		router.ServeHTTP(resp, req)
		
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), "Patients List")
	})
}

func TestNewPatient(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupPatientControllerTestDB(t)
	
	timestamp := time.Now().UnixNano()
	newPatientName := fmt.Sprintf("New Test Patient %d", timestamp)
	
	t.Cleanup(func() {
		db.Unscoped().Where("name = ?", newPatientName).Delete(&models.Patient{})
	})
	
	t.Run("NewPatient GET form as receptionist", func(t *testing.T) {
		router := gin.Default()
		router.LoadHTMLGlob("../templates/*")
		
		router.Use(func(c *gin.Context) {
			role, _ := c.Cookie("user")
			if role != "" { 
				c.Set("user", role) 
			}
			c.Next()
		})
		
		router.GET("/patients/new", NewPatient(db))
		
		resp := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/patients/new", nil)
		req.AddCookie(&http.Cookie{
			Name:  "user",
			Value: "receptionist",
		})
		router.ServeHTTP(resp, req)
		
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), "Add Patient")
	})
	
	t.Run("NewPatient GET as doctor redirects", func(t *testing.T) {
		router := gin.Default()
		router.LoadHTMLGlob("../templates/*")
		
		router.Use(func(c *gin.Context) {
			role, _ := c.Cookie("user")
			if role != "" { 
				c.Set("user", role) 
			}
			c.Next()
		})
		
		router.GET("/patients/new", NewPatient(db))
		
		resp := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/patients/new", nil)
		req.AddCookie(&http.Cookie{
			Name:  "user",
			Value: "doctor",
		})
		router.ServeHTTP(resp, req)
		
		assert.Equal(t, http.StatusSeeOther, resp.Code) 
		assert.Equal(t, "/patients", resp.Header().Get("Location"))
	})
	
	t.Run("NewPatient POST creates patient", func(t *testing.T) {
		router := gin.Default()
		router.LoadHTMLGlob("../templates/*")
		
		router.Use(func(c *gin.Context) {
			role, _ := c.Cookie("user")
			if role != "" { 
				c.Set("user", role) 
			}
			c.Next()
		})
		
		router.POST("/patients/new", NewPatient(db))
		
		form := url.Values{}
		form.Add("name", newPatientName)
		form.Add("age", "55")
		form.Add("detail", "New test detail")
		
		resp := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/patients/new", strings.NewReader(form.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "user",
			Value: "receptionist",
		})
		router.ServeHTTP(resp, req)
		
		assert.Equal(t, http.StatusSeeOther, resp.Code) 
		assert.Equal(t, "/patients", resp.Header().Get("Location"))
		
		var patient models.Patient
		db.Where("name = ?", newPatientName).First(&patient)
		assert.Equal(t, newPatientName, patient.Name)
		assert.Equal(t, 55, patient.Age)
	})
}

func TestEditPatient(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupPatientControllerTestDB(t)
	
	timestamp := time.Now().UnixNano()
	editPatientName := fmt.Sprintf("Edit Test Patient %d", timestamp)
	patient := models.Patient{
		Name:   editPatientName,
		Age:    40,
		Detail: "Original details",
	}
	db.Create(&patient)
	patientID := strconv.Itoa(int(patient.ID))
	
	t.Cleanup(func() {
		db.Unscoped().Where("name LIKE ?", editPatientName+"%").Delete(&models.Patient{})
	})
	
	t.Run("EditPatient GET shows edit form", func(t *testing.T) {
		router := gin.Default()
		router.LoadHTMLGlob("../templates/*")
		
		router.Use(func(c *gin.Context) {
			role, _ := c.Cookie("user")
			if role != "" { 
				c.Set("user", role) 
			}
			c.Next()
		})
		
		router.GET("/patients/edit/:id", EditPatient(db))
		
		resp := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/patients/edit/"+patientID, nil)
		req.AddCookie(&http.Cookie{
			Name:  "user",
			Value: "doctor",
		})
		router.ServeHTTP(resp, req)
		
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), "Edit Patient")
	})
	
	t.Run("EditPatient POST updates patient", func(t *testing.T) {
		router := gin.Default()
		router.LoadHTMLGlob("../templates/*")
		
		router.Use(func(c *gin.Context) {
			role, _ := c.Cookie("user")
			if role != "" { 
				c.Set("user", role) 
			}
			c.Next()
		})
		
		router.POST("/patients/edit/:id", EditPatient(db))
		
		updatedName := editPatientName + " Updated"
		
		form := url.Values{}
		form.Add("name", updatedName)
		form.Add("age", "60")
		form.Add("detail", "Updated details")
		
		resp := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/patients/edit/"+patientID, strings.NewReader(form.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "user",
			Value: "doctor",
		})
		router.ServeHTTP(resp, req)
		
		assert.Equal(t, http.StatusSeeOther, resp.Code) 
		assert.Equal(t, "/patients", resp.Header().Get("Location"))
		
		var updatedPatient models.Patient
		db.First(&updatedPatient, patientID)
		assert.Equal(t, updatedName, updatedPatient.Name)
		assert.Equal(t, 60, updatedPatient.Age)
		assert.Equal(t, "Updated details", updatedPatient.Detail)
	})
}

func TestDeletePatient(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupPatientControllerTestDB(t)
	
	timestamp := time.Now().UnixNano()
	deletePatientName := fmt.Sprintf("Delete Test Patient %d", timestamp)
	patient := models.Patient{
		Name:   deletePatientName,
		Age:    35,
		Detail: "To be deleted",
	}
	db.Create(&patient)
	patientID := strconv.Itoa(int(patient.ID))
	
	t.Cleanup(func() {
		db.Unscoped().Where("name = ?", deletePatientName).Delete(&models.Patient{})
	})
	
	t.Run("DeletePatient as receptionist deletes patient", func(t *testing.T) {
		router := gin.Default()
		
		router.Use(func(c *gin.Context) {
			role, _ := c.Cookie("user")
			if role != "" { 
				c.Set("user", role) 
			}
			c.Next()
		})
		
		router.POST("/patients/delete/:id", DeletePatient(db))
		
		resp := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/patients/delete/"+patientID, nil)
		req.AddCookie(&http.Cookie{
			Name:  "user",
			Value: "receptionist",
		})
		router.ServeHTTP(resp, req)
		
		assert.Equal(t, http.StatusSeeOther, resp.Code) 
		assert.Equal(t, "/patients", resp.Header().Get("Location"))
		
		var deletedPatient models.Patient
		result := db.First(&deletedPatient, patientID)
		assert.NotNil(t, result.Error)
		assert.Equal(t, "record not found", result.Error.Error())
	})
	
	t.Run("DeletePatient as doctor is prevented", func(t *testing.T) {
		
		anotherPatient := models.Patient{
			Name:   deletePatientName + " - Doctor Test",
			Age:    50,
			Detail: "Should not be deleted by doctor",
		}
		db.Create(&anotherPatient)
		anotherPatientID := strconv.Itoa(int(anotherPatient.ID))
		
		
		t.Cleanup(func() {
			db.Unscoped().Where("name = ?", deletePatientName+" - Doctor Test").Delete(&models.Patient{})
		})
		
		router := gin.Default()
		
		
		router.Use(func(c *gin.Context) {
			role, _ := c.Cookie("user")
			if role != "" { 
				c.Set("user", role) 
			}
			c.Next()
		})
		
		router.POST("/patients/delete/:id", DeletePatient(db))
		
		
		resp := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/patients/delete/"+anotherPatientID, nil)
		req.AddCookie(&http.Cookie{
			Name:  "user",
			Value: "doctor",
		})
		router.ServeHTTP(resp, req)
		
		assert.Equal(t, http.StatusSeeOther, resp.Code) 
		assert.Equal(t, "/patients", resp.Header().Get("Location"))
		
		
		var patient models.Patient
		result := db.First(&patient, anotherPatientID)
		assert.Nil(t, result.Error)
		assert.Equal(t, anotherPatient.Name, patient.Name)
	})
}