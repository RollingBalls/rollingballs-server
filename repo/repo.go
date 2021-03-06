package repo

import (
	"encoding/csv"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
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

	where := bson.M{
		"$geometry": bson.M{
			"type":        "Point",
			"coordinates": []float32{coordinatesAndDistance.Lon, coordinatesAndDistance.Lat},
		},
		"$maxDistance": coordinatesAndDistance.Distance,
	}

	if err := poiCollection.Find(bson.M{"position": bson.M{"$near": where}}).Select(bson.M{"puzzles": 1}).All(&points); err != nil {
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

func POICoordinates(id string) (types.Coordinates, string, error) {
	poiCollection := db.C("poi")
	var poi types.POI
	var coordinates types.Coordinates

	if err := poiCollection.FindId(bson.ObjectIdHex(id)).One(&poi); err != nil {
		return coordinates, "", err
	}

	return poi.Position, poi.Name, nil
}

func CheckDistance(userCoordinates, targetCoordinates types.Coordinates, threshold uint) (bool, error) {
	poiCollection := db.C("poi")
	points := []types.POI{}

	where := bson.M{
		"$geometry": bson.M{
			"type":        "Point",
			"coordinates": []float32{userCoordinates.Lon, userCoordinates.Lat},
		},
		"$maxDistance": threshold,
	}

	if err := poiCollection.Find(bson.M{"position": bson.M{"$near": where}}).All(&points); err != nil {
		return false, err
	}

	for _, point := range points {
		if point.Position == targetCoordinates {
			return true, nil
		}
	}

	return false, nil
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
			name, lat, lon := r[0], strings.Replace(r[1], ",", ".", -1), strings.Replace(r[2], ",", ".", -1)
			if lat == "" || lon == "" {
				continue
			}

			latf, err := strconv.ParseFloat(lat, 32)
			if err != nil {
				panic(err)
			}
			lonf, err := strconv.ParseFloat(lon, 32)
			if err != nil {
				panic(err)
			}

			var puzzles []string

			if r[5] != "" {
				puzzles = append(puzzles, r[5])
			}

			poiCollection.Upsert(bson.M{"name": name}, bson.M{
				"name":     name,
				"position": bson.M{"lat": float32(latf), "lng": float32(lonf)},
				"puzzles":  puzzles,
			})
		}

		if err = poiCollection.EnsureIndex(mgo.Index{Key: []string{"$2dsphere:position"}}); err != nil {
			log.Println(err)
		}

		log.Println("CSV refreshed")
		time.Sleep(5 * time.Minute)
	}
}
