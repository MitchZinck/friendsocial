package test_activity_locations

import (
	"friendsocial/activity_locations"
	"friendsocial/postgres"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServiceCreate(t *testing.T) {
	// Initialize the database connection
	postgres.InitDB()
	defer postgres.CloseDB()

	// Ensure the database connection is successful
	assert.NotNil(t, postgres.DB)

	service := activity_locations.NewService(postgres.DB)

	newLocation := activity_locations.ActivityLocation{
		Name:      "Park",
		Address:   "123 Green St",
		City:      "Greenwood",
		State:     "GW",
		ZipCode:   "12345",
		Country:   "Wonderland",
		Latitude:  1.2345,
		Longitude: 6.7890,
	}

	// Test the Create method
	createdLocation, err := service.Create(newLocation)
	assert.NoError(t, err)

	// Verify that the returned location matches the input
	assert.Equal(t, newLocation.Name, createdLocation.Name)
	assert.Equal(t, newLocation.Address, createdLocation.Address)
	assert.Equal(t, newLocation.City, createdLocation.City)
	assert.Equal(t, newLocation.State, createdLocation.State)
	assert.Equal(t, newLocation.ZipCode, createdLocation.ZipCode)
	assert.Equal(t, newLocation.Country, createdLocation.Country)
	assert.Equal(t, newLocation.Latitude, createdLocation.Latitude)
	assert.Equal(t, newLocation.Longitude, createdLocation.Longitude)
}
