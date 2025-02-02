package lib

import (
  "gorm.io/gorm"
  "github.com/lib/pq"
  "github.com/google/uuid"
)

type User struct {
	gorm.Model
	ID uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Email string `gorm:"unique"`
	Password string
}

type Car struct {
	gorm.Model
	Make         string
	ModelName        string
	Year         int
	EngineSize   float64
	FuelType     string
	Transmission string
	OwnerID      string
  ImageURLs    pq.StringArray `gorm:"type:text[];size:10"`
}

type CarReceived struct {
	Make         string
	ModelName        string
	Year         int
	EngineSize   float64
	FuelType     string
	Transmission string
	OwnerID      string
  ImageURLs    []byte
}

