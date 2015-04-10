package webserver

import (
	"log"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	cors "github.com/tommy351/gin-cors"
)

func Run(listen string) {
	gin.SetMode(gin.DebugMode)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(cors.Middleware(cors.Options{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"x-auth-token", "content-type"},
	}))

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
