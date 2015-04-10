package repo

import (
	"encoding/csv"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/RollingBalls/rollingballs-server/types"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var db *mgo.Database
var session *mgo.Session

func init() {
	var err error
	session, err = mgo.Dial("localhost")
	if err != nil {
		log.Fatal("Error connecting to mongo:", err)
	}

	session.SetMode(mgo.Monotonic, true)

	db = session.DB("rollingballs")

	go refreshOpenData()
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
		if len(point.Puzzles) > 0 {
			puzzleText := point.Puzzles[rand.Intn(len(point.Puzzles))]
			puzzle := map[string]string{"id": point.Id.Hex(), "text": puzzleText}
			puzzles = append(puzzles, puzzle)
		}
	}

	return puzzles, nil
}

func refreshOpenData() {

	const CSV_URL = "https://docs.google.com/spreadsheets/d/1T02iEmlUdnqEv2gpq_Y200TywLEMoxjI2D2EaBp1c9w/export?gid=0&format=csv"

	for {
		db = session.DB("rollingballs")
		time.Sleep(1 * time.Second)

		resp, err := http.Get(CSV_URL)
		if err != nil {
			log.Print("cannot fetch CSV", err)
			continue
		}

		csv := csv.NewReader(resp.Body)
		csv.Read() // skip first line

		records, err := csv.ReadAll()
		if err != nil {
			log.Print("cannot fetch CSV", err)
			continue
		}

		poiCollection := db.C("poi")
		for _, r := range records {
			// Skip records without lat/lon
			name, lat, lon := r[0], r[1], r[2]
			if lat == "" || lon == "" {
				continue
			}

			latf, err := strconv.ParseFloat(lat, 32)
			if err == nil {
				panic(err)
			}
			lonf, err := strconv.ParseFloat(lon, 32)
			if err == nil {
				panic(err)
			}

			poiCollection.Upsert(bson.M{"name": name}, bson.M{"name": name, "lat": latf, "lon": lonf})
		}

		log.Println("CSV refreshed")
		time.Sleep(5 * time.Minute)
	}
}
