package engine

import (
	"github.com/RollingBalls/rollingballs-server/types"
	"time"
)

type game struct {
	startTime time.Time
	target    types.Coordinates
}

type Engine struct {
	games map[string]game
}

func New() (engine *Engine) {
	engine = new(Engine)
	engine.games = make(map[string]game)

	return
}

func (engine *Engine) StartGame(userToken string, targetCoordinates types.Coordinates) uint {
	engine.games[userToken] = game{time.Now(), targetCoordinates}
	return 3600
}
