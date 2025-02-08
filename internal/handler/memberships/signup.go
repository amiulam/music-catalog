package memberships

import (
	"net/http"

	"github.com/amiulam/music-catalog/internal/models/memberships"
	"github.com/gin-gonic/gin"
)

func (h *Handler) SignUp(c *gin.Context) {
	var request memberships.SignUpRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := h.membershipSvc.SignUp(request)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user created successfully",
	})
}
