package httpapi

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jesrnsu/jesrnsu.github.io/internal/rounds"
)

type roundTotalResponse struct {
	TotalStrokes       int `json:"total_strokes"`
	TotalPar           int `json:"total_par"`
	ScoreRelativeToPar int `json:"score_relative_to_par"`
}

type roundParResponse struct {
	TotalPar int `json:"total_par"`
}

func getUserRounds(service *rounds.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Query("username")
		if strings.TrimSpace(username) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
			return
		}
		if service == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "round storage is unavailable"})
			return
		}
		userRounds, err := service.GetUserRounds(c.Request.Context(), username)
		if err != nil {
			_ = c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get rounds"})
			return
		}
		c.JSON(http.StatusOK, userRounds)
	}
}

func saveUserRound(service *rounds.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		if service == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "round storage is unavailable"})
			return
		}
		var round rounds.Round
		if err := c.ShouldBindJSON(&round); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "request body must match the round JSON structure"})
			return
		}
		if err := rounds.ValidateRound(round); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		id, err := service.SaveRound(c.Request.Context(), round)
		if err != nil {
			_ = c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save round"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"id": id, "message": "round saved"})
	}
}

func calculateRoundTotal(c *gin.Context) {
	var round rounds.Round
	if err := c.ShouldBindJSON(&round); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "request body must match the round JSON structure"})
		return
	}

	// Validate the completed round.
	if err := rounds.ValidateRound(round); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Translate the result into an HTTP response.
	c.JSON(http.StatusOK, roundTotalResponse{
		TotalStrokes:       rounds.TotalStrokes(round),
		TotalPar:           rounds.TotalPar(round),
		ScoreRelativeToPar: rounds.ScoreRelativeToPar(round),
	})
}

func calculateCoursePar(c *gin.Context) {
	var round rounds.Round
	if err := c.ShouldBindJSON(&round); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "request body must match the round JSON structure"})
		return
	}

	if err := rounds.ValidateRound(round); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	total := rounds.TotalPar(round)

	c.JSON(http.StatusOK, roundParResponse{TotalPar: total})
}
