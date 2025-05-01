package controllers

import (
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "github.com/mani2002/golang-clinic-portal/models"
    "golang.org/x/crypto/bcrypt"
    "net/http"
)

func ShowLogin(c *gin.Context) {
    c.HTML(http.StatusOK, "login.html", gin.H{"error": ""})
}

func PostLogin(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        email := c.PostForm("email")
        password := c.PostForm("password")
        role := c.PostForm("role")
        var user models.User
        if err := db.Where("email = ? AND role = ?", email, role).First(&user).Error; err != nil {
            c.HTML(http.StatusOK, "login.html", gin.H{"error": "Invalid credentials"})
            return
        }
        if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
            c.HTML(http.StatusOK, "login.html", gin.H{"error": "Invalid credentials"})
            return
        }
        c.SetCookie("user", user.Role, 3600, "/", "", false, true)
        c.Redirect(http.StatusSeeOther, "/patients")
    }
}