package engine

import (
	"errors"
	"github.com/RollingBalls/rollingballs-server/types"
	"math/rand"
	"time"
)

type game struct {
	startTime time.Time
	target    types.Coordinates
}

type Engine struct {
	games map[string]*game
}

func New() (engine *Engine) {
	engine = new(Engine)
	engine.games = make(map[string]*game)

	return
}

func (engine *Engine) Target(userToken string) (types.Coordinates, error) {
	if engine.games[userToken] == nil {
		return types.Coordinates{}, errors.New("Not Found")
	}

	return engine.games[userToken].target, nil
}

func (engine *Engine) StartGame(userToken string, targetCoordinates types.Coordinates) uint {
	engine.games[userToken] = &game{time.Now(), targetCoordinates}

	return uint(rand.Intn(30) + 30)
}

func (engine *Engine) EndGame(userToken string) error {
	if engine.games[userToken] == nil {
		return errors.New("Not Found")
	}

	delete(engine.games, userToken)
	return nil
}
