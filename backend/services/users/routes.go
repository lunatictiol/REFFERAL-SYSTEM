package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/lunatictiol/referal-system/config"
	"github.com/lunatictiol/referal-system/models"
	"github.com/lunatictiol/referal-system/services/auth"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Handler struct {
	store UserStore
}

var Validator = validator.New()

func NewHandler(store UserStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) RegisterRoutes(route *gin.RouterGroup) {
	route.GET("/ping", h.ping)
	route.POST("/user/register", h.registerUser)
	route.POST("/user/login", h.loginUser)

}

func (h *Handler) ping(c *gin.Context) {

	c.JSON(200, gin.H{
		"message": "pong",
	})
	h.store.Ping()

}

func (h *Handler) registerUser(c *gin.Context) {
	var payload models.RegisterUserPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "cant parse json",
		})
		return
	}
	if err := Validator.Struct(payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "invalid data",
		})
		return
	}
	existingUser, err := h.store.CheckIfEmailExisits(payload.Email)
	if err != nil {
		config.LogFatal("Error finding user")
	}
	if existingUser.Id != bson.NilObjectID {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "user already exists",
		})
		return

	}
	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "error creating the user",
		})
		return
	}

	user := models.RegisterUserPayload{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword,
	}

	uid, err := h.store.RegisterUser(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "error creating user",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":  true,
		"user_id": uid,
	})

}
func (h *Handler) loginUser(c *gin.Context) {
	var payload models.LoginUserPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "cant parse json",
		})
		return
	}
	if err := Validator.Struct(payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "invalid data",
		})
		return
	}
	existingUser, err := h.store.CheckIfEmailExisits(payload.Email)
	if err != nil {
		config.LogFatal("Error finding user")
	}
	if existingUser.Id == bson.NilObjectID {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "user does'nt exists",
		})
		return

	}

	if !auth.ValidatePassword(payload.Password, existingUser.Password) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "error getting data of the user",
		})
		return
	}
	token, err := auth.GenerateToken(existingUser.Id.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "error getting data of the user",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "login successfull",
		"token":   token,
	})

}
