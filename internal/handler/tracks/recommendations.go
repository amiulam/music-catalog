package tracks

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetRecommendation(c *gin.Context) {
	ctx := c.Request.Context()

	userID := c.GetUint("userID")
	trackID := c.Query("trackID")
	limit, err := strconv.Atoi(c.Query("limit"))

	if err != nil {
		limit = 10
	}

	response, err := h.trackSvc.GetRecommendation(ctx, limit, trackID, userID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}
