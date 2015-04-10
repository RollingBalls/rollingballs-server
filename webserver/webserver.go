package webserver

import (
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"log"
)

func Run(listen string) {
	gin.SetMode(gin.DebugMode)

	r := gin.New()
	r.Use(gin.Recovery())

	apiEndpoints := r.Group("/api/v1")
	{
		apiEndpoints.GET("/puzzles", puzzlesByPositionAndDistance)
		apiEndpoints.GET("/remaining-time", remainingTime)
		apiEndpoints.PUT("/start", startGame)
		apiEndpoints.PUT("/abort", abortGame)
		apiEndpoints.PUT("/finish", finishGame)
	}

	r.Use(static.Serve("/", static.LocalFile("static", false)))

	if error := r.Run(listen); error != nil {
		log.Fatal(error)
	}
}
