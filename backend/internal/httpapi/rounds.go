package httpapi

import (
	"errors"
	"net/http"
	"strconv"

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

func getUserRound(service *rounds.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil || id <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "round id must be a positive integer"})
			return
		}
		if service == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "round storage is unavailable"})
			return
		}

		round, err := service.GetRound(c.Request.Context(), id)
		if errors.Is(err, rounds.ErrRoundNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "round not found"})
			return
		}
		if err != nil {
			_ = c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get round"})
			return
		}
		c.JSON(http.StatusOK, round)
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
