package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/lunatictiol/referal-system/models"
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
	route.GET("/user/register", h.registerUser)

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
			"error": err.Error(),
		})
		return
	}
	if err := Validator.Struct(payload); err != nil {
		err := err.(validator.ValidationErrors)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

}
