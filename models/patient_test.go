package models

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"fmt"
	"time"
)

func setupPatientTestDB(t *testing.T) *gorm.DB {
	dsn := "host=localhost user=postgres password=postgres dbname=clinicdb port=5432 sslmode=disable"
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	
	return db
}

func TestPatientModel(t *testing.T) {
	db := setupPatientTestDB(t)
	
	timestamp := time.Now().UnixNano()
	testName := fmt.Sprintf("Test Patient %d", timestamp)
	
	t.Run("Create patient", func(t *testing.T) {
		patient := Patient{
			Name:   testName,
			Age:    35,
			Detail: "Test detail only",
		}
		
		result := db.Create(&patient)
		
		t.Cleanup(func() {
			db.Unscoped().Where("name = ?", testName).Delete(&Patient{})
		})
		
		assert.Nil(t, result.Error)
		assert.NotZero(t, patient.ID)
		assert.Equal(t, testName, patient.Name)
		assert.Equal(t, 35, patient.Age)
	})
	
	t.Run("Update patient", func(t *testing.T) {

		updateTestName := fmt.Sprintf("Update Test Patient %d", timestamp)
		patient := Patient{
			Name:   updateTestName,
			Age:    42,
			Detail: "Chronic condition",
		}
		db.Create(&patient)
		
		t.Cleanup(func() {
			db.Unscoped().Where("name LIKE ?", updateTestName+"%").Delete(&Patient{})
		})
		
		patientID := patient.ID
		patient.Name = updateTestName + " Updated"
		patient.Age = 43
		db.Save(&patient)
		
		var updatedPatient Patient
		db.First(&updatedPatient, patientID)
		
		assert.Equal(t, updateTestName+" Updated", updatedPatient.Name)
		assert.Equal(t, 43, updatedPatient.Age)
	})
	
	t.Run("Delete patient", func(t *testing.T) {

		deleteTestName := fmt.Sprintf("Delete Test Patient %d", timestamp)
		patient := Patient{
			Name:   deleteTestName,
			Age:    50,
			Detail: "To be deleted",
		}
		db.Create(&patient)
		patientID := patient.ID
		
		t.Cleanup(func() {
			db.Unscoped().Where("name = ?", deleteTestName).Delete(&Patient{})
		})
		
		db.Delete(&Patient{}, patientID)
		
		var deletedPatient Patient
		result := db.First(&deletedPatient, patientID)
		
		assert.NotNil(t, result.Error)
		assert.Equal(t, "record not found", result.Error.Error())
	})
}