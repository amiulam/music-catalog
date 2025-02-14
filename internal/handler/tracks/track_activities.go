package tracks

import (
	"net/http"

	trackactivities "github.com/amiulam/music-catalog/internal/models/track_activities"
	"github.com/gin-gonic/gin"
)

func (h *Handler) UpsertTrackActivities(c *gin.Context) {
	ctx := c.Request.Context()

	var request trackactivities.TrackActivityRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userID := c.GetUint("userID")

	err := h.trackSvc.UpsertUserActivities(ctx, userID, request)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "activity saved!",
	})
}
