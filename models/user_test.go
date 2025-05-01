package models

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	dsn := "host=localhost user=postgres password=postgres dbname=clinicdb port=5432 sslmode=disable"
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	
	return db
}

func TestUserModel(t *testing.T) {
	db := setupTestDB(t)
	
	testEmail := "test_user_" + randomString(8) + "@example.com"
	
	testUser := User{
		Email:    testEmail,
		Password: "testpassword",
		Role:     "doctor",
	}
	
	t.Cleanup(func() {
		db.Unscoped().Where("email = ?", testEmail).Delete(&User{})
	})
	
	t.Run("Create user", func(t *testing.T) {
		result := db.Create(&testUser)
		
		assert.Nil(t, result.Error)
		assert.NotZero(t, testUser.ID)
		assert.Equal(t, testEmail, testUser.Email)
		assert.Equal(t, "doctor", testUser.Role)
	})
	
	t.Run("Find user by email and role", func(t *testing.T) {

		var foundUser User
		result := db.Where("email = ? AND role = ?", testEmail, "doctor").First(&foundUser)
		
		assert.Nil(t, result.Error)
		assert.Equal(t, testEmail, foundUser.Email)
		assert.Equal(t, "doctor", foundUser.Role)
	})
	
	t.Run("Email uniqueness constraint", func(t *testing.T) {

		duplicateUser := User{
			Email:    testEmail, 
			Password: "anotherpassword",
			Role:     "receptionist",
		}
		
		result := db.Create(&duplicateUser)
		
		assert.NotNil(t, result.Error)
	})
}

func randomString(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[i%len(chars)]
	}
	return string(result)
}