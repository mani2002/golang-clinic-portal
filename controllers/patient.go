package controllers

import (
	"github.com/mani2002/golang-clinic-portal/models"
	"gorm.io/gorm"
	"net/http"
	"github.com/gin-gonic/gin"
	"strconv"
)

func ListPatients(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var patients []models.Patient
		db.Order("created_at desc").Find(&patients)
		
		
		role, exists := c.Get("user")
		if !exists {
			role = ""
		}
		
		c.HTML(http.StatusOK, "patients.html", gin.H{"patients": patients, "role": role})
	}
}

func NewPatient(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		
		role, exists := c.Get("user")
		if !exists {
			role = ""
		}
		
		
		if role != "receptionist" { 
			c.Redirect(http.StatusSeeOther, "/patients")
			return 
		}
		
		if c.Request.Method == http.MethodPost {
			name := c.PostForm("name")
			age, _ := strconv.Atoi(c.PostForm("age"))
			detail := c.PostForm("detail")
			db.Create(&models.Patient{Name: name, Age: age, Detail: detail})
			c.Redirect(http.StatusSeeOther, "/patients")
			return
		}
		
		c.HTML(http.StatusOK, "edit_patient.html", gin.H{"patient": models.Patient{}, "role": role})
	}
}

func EditPatient(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var patient models.Patient
		if db.First(&patient, id).Error != nil {
			c.String(404, "Not found")
			return
		}
		
		
		role, exists := c.Get("user")
		if !exists {
			role = ""
		}
		
		if c.Request.Method == http.MethodPost {
			patient.Name = c.PostForm("name")
			patient.Age, _ = strconv.Atoi(c.PostForm("age"))
			patient.Detail = c.PostForm("detail")
			db.Save(&patient)
			c.Redirect(http.StatusSeeOther, "/patients")
			return
		}
		
		c.HTML(http.StatusOK, "edit_patient.html", gin.H{"patient": patient, "role": role})
	}
}

func DeletePatient(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		
		role, exists := c.Get("user")
		if !exists || role != "receptionist" { 
			c.Redirect(http.StatusSeeOther, "/patients")
			return 
		}
		
		id := c.Param("id")
		db.Delete(&models.Patient{}, id)
		c.Redirect(http.StatusSeeOther, "/patients")
	}
}