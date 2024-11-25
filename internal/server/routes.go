package server

import (
	"net/http"
	"referalsystem/internal/types"
	"referalsystem/internal/utils"

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
