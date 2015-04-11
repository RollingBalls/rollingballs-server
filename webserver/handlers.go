package webserver

import (
	"net/http"
	"strconv"

	"github.com/RollingBalls/rollingballs-server/engine"
	"github.com/RollingBalls/rollingballs-server/repo"
	"github.com/RollingBalls/rollingballs-server/types"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

const WinThreshold = 50

type JSONObject map[string]interface{}

func puzzlesByPositionAndDistance(c *gin.Context) {
	coordinates := types.CoordinatesAndDistance{Coordinates: &types.Coordinates{}}

	if lat, err := strconv.ParseFloat(c.Request.URL.Query().Get("lat"), 32); err != nil {
		c.Fail(http.StatusBadRequest, err)
		return
	} else {
		coordinates.Lat = float32(lat)
	}
	if lon, err := strconv.ParseFloat(c.Request.URL.Query().Get("lon"), 32); err != nil {
		c.Fail(http.StatusBadRequest, err)
		return
	} else {
		coordinates.Lon = float32(lon)
	}
	if distance, err := strconv.ParseUint(c.Request.URL.Query().Get("distance"), 10, 32); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	} else {
		coordinates.Distance = uint(distance)
	}

	if coordinates.Valid() {
		if puzzles, error := repo.Puzzles(coordinates); error != nil {
			c.AbortWithStatus(http.StatusBadRequest)
		} else {
			c.JSON(http.StatusOK, puzzles)
		}
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func remainingTime(c *gin.Context) {
	// TODO: implement
	c.JSON(http.StatusOK, JSONObject{})
}

func startGame(c *gin.Context) {
	var submission GuessingSubmission
	c.BindWith(&submission, binding.JSON)
	engine := c.MustGet("engine").(*engine.Engine)

	if submission.Id != "" {
		if coordinates, name, err := repo.POICoordinates(submission.Id); err != nil {
			c.Fail(http.StatusBadRequest, err)
		} else {
			seconds := engine.StartGame(c.Request.Header.Get("X-Auth-Token"), coordinates, name)
			c.JSON(http.StatusCreated, JSONObject{"seconds": seconds})
		}
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func abortGame(c *gin.Context) {
	c.JSON(http.StatusNoContent, nil)
}

func finishGame(c *gin.Context) {
	var submission PositionSubmission
	c.BindWith(&submission, binding.JSON)
	engine := c.MustGet("engine").(*engine.Engine)

	coordinates := types.Coordinates{submission.Lat, submission.Lon}

	if coordinates.Valid() {
		if target, name, err := engine.Target(c.Request.Header.Get("X-Auth-Token")); err != nil {
			c.Fail(http.StatusNotFound, err)
		} else {
			win, _ := repo.CheckDistance(coordinates, target, WinThreshold)
			if err := engine.EndGame(c.Request.Header.Get("X-Auth-Token")); err != nil {
				c.Fail(http.StatusNotFound, err)
			} else {
				c.JSON(http.StatusOK, JSONObject{"win": win, "poi": name})
			}
		}
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}
