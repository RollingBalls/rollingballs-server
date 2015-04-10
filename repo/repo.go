package repo

import (
	"github.com/RollingBalls/rollingballs-server/types"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"math/rand"
)

var db *mgo.Database

func init() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatal("Error connecting to mongo:", err)
	}

	session.SetMode(mgo.Monotonic, true)

	db = session.DB("rollingballs")
}

func Puzzles(coordinatesAndDistance types.CoordinatesAndDistance) ([]map[string]string, error) {
	poiCollection := db.C("poi")
	points := []types.POI{}
	puzzles := []map[string]string{}

	//where := bson.M{"$near": []float32{coordinatesAndDistance.Lat, coordinatesAndDistance.Lon}, "$maxDistance": coordinatesAndDistance.Distance}

	if err := poiCollection.Find(nil).Select(bson.M{"puzzles": 1}).All(&points); err != nil {
		log.Println(err)
		return puzzles, err
	}

	for _, point := range points {
		puzzleText := point.Puzzles[rand.Intn(len(point.Puzzles))]
		puzzle := map[string]string{"id": point.Id.Hex(), "text": puzzleText}
		puzzles = append(puzzles, puzzle)
	}

	return puzzles, nil
}
