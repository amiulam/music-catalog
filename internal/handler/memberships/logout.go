package memberships

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Logout(c *gin.Context) {
	c.Header("Authorization", "")

	c.JSON(http.StatusOK, gin.H{
		"message": "logout success",
	})
}
