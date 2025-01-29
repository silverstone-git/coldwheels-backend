package main

import (
	"coldwheels/middleware"
	"coldwheels/models"
	imageUpload "coldwheels/repository"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/argon2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type CarRequest struct {
	Make         string   `json:"make"`
	ModelName        string   `json:"modelName"`
	Year         int      `json:"year"`
	EngineSize   float64  `json:"engineSize"`
	FuelType     string   `json:"fuelType"`
	Transmission string   `json:"transmission"`
	ImageURLs    []string `json:"imageUrls"`
}

var JwtSecret []byte
var db *gorm.DB

func main() {

  err := godotenv.Load()
  if err != nil {
      fmt.Println("Error loading .env file")
      return
  }
  JwtSecret = []byte(os.Getenv("JWT_SECRET"))

  host := os.Getenv("DB_HOST")
  if host == "" {
        panic("DATABASE_URL environment variable is not set")
  }

  port := "5432"
  user := os.Getenv("DB_USER")
  password := os.Getenv("DB_PASSWORD")
  dbname := os.Getenv("DB_NAME")
  sslmode := os.Getenv("DB_SSL")


  dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        host, port, user, password, dbname, sslmode)


  db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}


	// Migrate the schema
	// db.AutoMigrate(&models.User{}, &models.Car{})

	// Initialize Gin
	r := gin.Default()

	// Auth routes
	r.POST("/signup", signup)
	r.POST("/login", login)

	// Protected routes
	auth := r.Group("/")

  auth.Use(middleware.AuthMiddleware(JwtSecret))
	{
		auth.GET("/api/cars", getCars)
		auth.POST("/api/cars", createCar)
		auth.PUT("/api/cars/:id", updateCar)
		auth.DELETE("/api/cars/:id", deleteCar)
		auth.POST("/api/cars/upload-images", imageUpload.UploadImagesHandler)
	}

	// Start server
	r.Run(":" + os.Getenv("PORT"))
}

// Auth Handlers
func signup(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

  user.ID = uuid.New()

	hashedPassword, _ := HashPassword(user.Password)
	user.Password = hashedPassword

  if db == nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection is not initialized"})
    return
}

	if result := db.Create(&user); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func login(c *gin.Context) {
	var credentials struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if result := db.Where("email = ?", credentials.Email).First(&user); result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !CheckPasswordHash(credentials.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(JwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// Car Handlers
func getCars(c *gin.Context) {
	var cars []models.Car
	userID := c.MustGet("userID").(string)
	db.Where("owner_id = ?", userID).Find(&cars)
	c.JSON(http.StatusOK, cars)
}

func createCar(c *gin.Context) {
	var req CarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate image URLs
	if len(req.ImageURLs) > 10 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Maximum 10 image URLs allowed"})
		return
	}

	userID := c.MustGet("userID").(string)

  // Check the number of cars owned by the user
	var carCount int64
	db.Model(&models.Car{}).Where("owner_id = ?", userID).Count(&carCount)
	if carCount >= 20 {
		c.JSON(http.StatusForbidden, gin.H{"error": "You cannot own more than 20 cars"})
		return
	}

	car := models.Car{
		Make:         req.Make,
		ModelName:        req.ModelName,
		Year:         req.Year,
		EngineSize:   req.EngineSize,
		FuelType:     req.FuelType,
		Transmission: req.Transmission,
		OwnerID:      userID,
		ImageURLs:    req.ImageURLs,
	}

	if result := db.Create(&car); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, car)
}

func updateCar(c *gin.Context) {
	id := c.Param("id")
	var req CarRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.ImageURLs) > 10 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Maximum 10 image URLs allowed"})
		return
	}

	var car models.Car
	if result := db.First(&car, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Car not found"})
		return
	}

	if car.OwnerID != c.MustGet("userID").(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	// Update fields
	car.Make = req.Make
	car.ModelName = req.ModelName
	car.Year = req.Year
	car.EngineSize = req.EngineSize
	car.FuelType = req.FuelType
	car.Transmission = req.Transmission
	car.ImageURLs = req.ImageURLs

	db.Save(&car)
	c.JSON(http.StatusOK, car)
}


func deleteCar(c *gin.Context) {
	id := c.Param("id")
	var car models.Car

	if result := db.First(&car, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Car not found"})
		return
	}

	if car.OwnerID != c.MustGet("userID").(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	db.Delete(&car)
	c.JSON(http.StatusOK, gin.H{"message": "Car deleted"})
}

// Helpers
func HashPassword(password string) (string, error) {

  salt := make([]byte, 16)

  if _, err := rand.Read(salt); err != nil {
        return "", fmt.Errorf("lafde in making salt")
  }

  hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

  if hash == nil || len(hash) == 0 {
   err := errors.New("something went wrong")
   return "", err
  }

	// Convert the hash to a base64-encoded string
	saltString := base64.RawStdEncoding.EncodeToString(salt)
	hashString := base64.RawStdEncoding.EncodeToString(hash)

  output := fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s", 64, 1, 4, saltString, hashString)

  fmt.Println("hash hash ke")
  fmt.Println(output)
  return output, nil
}

func CheckPasswordHash(password, hashString string) bool {

  parts := strings.Split(hashString, "$")
  if len(parts) != 6 {
        return false
  }

  salt, _ := base64.RawStdEncoding.DecodeString(parts[4])
  hashBytes, _ := base64.RawStdEncoding.DecodeString(parts[5])

  // Generate a new hash with the provided password and the stored salt
  newHash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, uint32(len(hashBytes)))

  // Compare the hashes
  return string(newHash) == string(hashBytes)
}
