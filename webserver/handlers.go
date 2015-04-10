package webserver

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"strconv"
)

// TODO: move these structs away
type coordinates struct {
	Lat float32
	Lon float32
}

func (c coordinates) Valid() bool {
	return c.Lat >= -90.0 && c.Lat <= 90.0 && c.Lon >= -180.0 && c.Lon <= 180.0
}

type coordinatesAndDistance struct {
	*coordinates
	Distance uint
}

func (c coordinatesAndDistance) Valid() bool {
	return c.coordinates.Valid() && c.Lon <= 180.0
}

type JSONObject map[string]interface{}

func puzzlesByPositionAndDistance(c *gin.Context) {
	coordinates := coordinatesAndDistance{coordinates: &coordinates{}}

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
		// TODO: implement
		c.JSON(http.StatusOK, make([]JSONObject, 0))
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

	if submission.Id > 0 {
		// TODO: implement
		c.JSON(http.StatusCreated, JSONObject{})
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

	coordinates := coordinates{submission.Lat, submission.Lon}

	if coordinates.Valid() {
		// TODO: implement
		c.JSON(http.StatusOK, JSONObject{})
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}
