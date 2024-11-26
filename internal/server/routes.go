package server

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"referalsystem/internal/types"
	"referalsystem/internal/utils"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	r.GET("/", s.HelloWorldHandler)
	r.GET("/health", s.healthHandler)
	r.POST("/register", s.RegisterUser)
	r.POST("/login", s.LoginUser)
	r.POST("/generateReferal", s.ReferHandler)

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}

func (s *Server) RegisterUser(c *gin.Context) {
	var payload types.RegisterUserPayload
	refered := c.Query("refered")

	if refered == "true" {
		//todo
	}

	// get body
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": "invalid data",
		})
		return
	}

	// check if user exists
	_, err = s.db.FindUserByEmail(payload.Email)
	if err == nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": "cannot create user at the momentssss",
		})
		return
	}

	//hash password
	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "cannot create user at the moment",
		})
		return
	}

	//create user
	uID, err := s.db.CreateUser(types.User{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: hashedPassword,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "cannot create user at the moment",
		})
		return
	}

	c.JSON(http.StatusOK, map[string]string{
		"message": "Registeration successful",
		"user_id": uID,
	})
}

func (s *Server) LoginUser(c *gin.Context) {
	var payload types.LoginUserPayload

	// get body
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": "invalid data",
		})
		return
	}

	// check if user exists
	u, err := s.db.FindUserByEmail(payload.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": "cannot find the user with that email",
		})
		return
	}

	//validate  password
	validated := utils.ValidatePassword(payload.Password, u.Password)
	if !validated {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": "incorect credentials",
		})
		return
	}

	//create user

	c.JSON(http.StatusOK, map[string]string{
		"message": "Login successful",
		"user_id": u.ID,
	})
}

func (s *Server) ReferHandler(c *gin.Context) {
	var req types.ReferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request. 'service_id' is required."})
		return
	}
	// Channels for Goroutine communication
	resultChan := make(chan string)
	errorChan := make(chan error)

	// Start Goroutine to generate the referral code
	go GenerateReferralCode(req.ServiceID, req.UserID, resultChan, errorChan, s)

	select {
	case code := <-resultChan:
		c.JSON(http.StatusOK, gin.H{
			"message":       "Referral code generated successfully",
			"referral_code": code,
		})
		return
	case err := <-errorChan:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate referral code",
			"details": err.Error(),
		})
		return
	}
}

func GenerateReferralCode(serviceID, userId string, resultChan chan<- string, errorChan chan<- error, s *Server) {
	const randomPartLength = 8 // Length of the random string

	for {
		// Generate a random string
		randomBytes := make([]byte, randomPartLength)
		_, err := rand.Read(randomBytes)
		if err != nil {
			errorChan <- fmt.Errorf("failed to generate random string: %w", err)
			return
		}
		randomPart := hex.EncodeToString(randomBytes)

		// Add a timestamp for additional uniqueness
		timestamp := time.Now().UnixNano()

		// Combine service ID, random part, and timestamp
		referralCode := fmt.Sprintf("%s-%s-%d", serviceID, randomPart, timestamp)

		// Check if the code already exists in the database

		// Ensure uniqueness by checking against the global set

		exists, err := s.db.CheckReferralCodeExists(referralCode)
		if err != nil {
			errorChan <- fmt.Errorf("failed to check referral code existence: %w", err)
			return
		}
		if !exists {

			err = s.db.InsertReferralCodeExists(referralCode, userId)
			if err != nil {
				errorChan <- fmt.Errorf("failed to insert referral code: %w", err)
				return
			}
			resultChan <- referralCode
			return
		}

	}
}
