// Package httpapi translates between HTTP requests and application services.
package httpapi

import (
	"github.com/gin-gonic/gin"
	"github.com/jesrnsu/jesrnsu.github.io/internal/rounds"
)

func NewRouter(roundService *rounds.Service) *gin.Engine {
	router := gin.Default()

	router.Use(func(c *gin.Context) {
    c.Header("Access-Control-Allow-Origin", "http://localhost:4321")
    c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
    c.Header("Access-Control-Allow-Headers", "Content-Type")

    if c.Request.Method == "OPTIONS" {
        c.AbortWithStatus(204)
        return
    }

    c.Next()
	})

	router.POST("/api/v1/rounds/total", calculateRoundTotal)
	router.POST("/api/v1/rounds/course_par", calculateCoursePar)
	router.POST("/api/v1/rounds", saveUserRound(roundService))
	router.GET("/api/v1/rounds/:id", getUserRound(roundService))
	return router
}
