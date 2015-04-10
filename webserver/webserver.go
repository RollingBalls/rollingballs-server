package webserver

import (
	"github.com/RollingBalls/rollingballs-server/engine"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	cors "github.com/tommy351/gin-cors"
	"log"
)

func Run(listen string, engine *engine.Engine) {
	gin.SetMode(gin.DebugMode)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(RollingBalls(engine))
	r.Use(cors.Middleware(cors.Options{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"x-auth-token", "content-type"},
	}))

	apiEndpoints := r.Group("/api/v1")
	{
		apiEndpoints.GET("/puzzles", puzzlesByPositionAndDistance)
		authorizedEndpoints := apiEndpoints.Group("/", CheckAccessToken())
		{
			authorizedEndpoints.GET("/remaining-time", remainingTime, CheckAccessToken())
			authorizedEndpoints.PUT("/start", startGame, CheckAccessToken())
			authorizedEndpoints.PUT("/abort", abortGame, CheckAccessToken())
			authorizedEndpoints.PUT("/finish", finishGame, CheckAccessToken())
		}
	}

	r.Use(static.Serve("/", static.LocalFile("static", false)))

	if error := r.Run(listen); error != nil {
		log.Fatal(error)
	}
}
