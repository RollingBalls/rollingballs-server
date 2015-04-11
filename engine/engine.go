package engine

import (
	"errors"
	"math/rand"
	"time"

	"github.com/RollingBalls/rollingballs-server/types"
)

type game struct {
	startTime  time.Time
	target     types.Coordinates
	targetName string
}

type Engine struct {
	games map[string]*game
}

func New() (engine *Engine) {
	engine = new(Engine)
	engine.games = make(map[string]*game)

	return
}

func (engine *Engine) Target(userToken string) (types.Coordinates, string, error) {
	if engine.games[userToken] == nil {
		return types.Coordinates{}, "", errors.New("Not Found")
	}

	return engine.games[userToken].target, engine.games[userToken].targetName, nil
}

func (engine *Engine) StartGame(userToken string, targetCoordinates types.Coordinates, targetName string) uint {
	engine.games[userToken] = &game{time.Now(), targetCoordinates, targetName}

	return uint(rand.Intn(30) + 30)
}

func (engine *Engine) EndGame(userToken string) error {
	if engine.games[userToken] == nil {
		return errors.New("Not Found")
	}

	delete(engine.games, userToken)
	return nil
}
