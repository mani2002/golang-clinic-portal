package controllers

import (
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "github.com/mani2002/golang-clinic-portal/middleware"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
    r.GET("/login", ShowLogin)
    r.POST("/login", PostLogin(db))

    r.Use(func(c *gin.Context) {
        role, _ := c.Cookie("user")
        c.Set("user", role)
        c.Next()
    })

    r.GET("/patients", ListPatients(db))
    r.GET("/patients/new", middleware.RequireAuth("receptionist"), NewPatient(db))
    r.POST("/patients/new", middleware.RequireAuth("receptionist"), NewPatient(db))
	r.GET("/patients/edit/:id", EditPatient(db))
	r.POST("/patients/edit/:id", EditPatient(db))
    r.POST("/patients/delete/:id", middleware.RequireAuth("receptionist"), DeletePatient(db))
}