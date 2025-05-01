package main

import (
    "github.com/gin-gonic/gin"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "os"
    "github.com/joho/godotenv"
    "github.com/mani2002/golang-clinic-portal/models"
    "github.com/mani2002/golang-clinic-portal/controllers"
   // "github.com/mani2002/golang-clinic-portal/middleware"
)

func setupDatabase() *gorm.DB {
    dsn := "host=" + os.Getenv("DB_HOST") + " user=" + os.Getenv("DB_USER") +
        " password=" + os.Getenv("DB_PASSWORD") + " dbname=" + os.Getenv("DB_NAME") +
        " port=" + os.Getenv("DB_PORT") + " sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("Failed to connect database")
    }
    db.AutoMigrate(&models.User{}, &models.Patient{})
    return db
}

func main() {
    godotenv.Load()
    db := setupDatabase()
    r := gin.Default()
    r.Static("/static", "./static")
    r.LoadHTMLGlob("templates/*")
    controllers.SetupRoutes(r, db)
    r.Run(":8080")
}