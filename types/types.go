package types

import (
	"gopkg.in/mgo.v2/bson"
)

type Coordinates struct {
	Lat float32
	Lon float32
}

func (c Coordinates) Valid() bool {
	return c.Lat >= -90.0 && c.Lat <= 90.0 && c.Lon >= -180.0 && c.Lon <= 180.0
}

type CoordinatesAndDistance struct {
	*Coordinates
	Distance uint
}

func (c CoordinatesAndDistance) Valid() bool {
	return c.Coordinates.Valid() && c.Lon <= 180.0
}

type POI struct {
	Id      bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name    string        `json:"name,omitempty" bson:"name,omitempty`
	Lat     float64       `json:"lat,omitempty" bson:"lat,omitempty`
	Lon     float64       `json:"lon,omitempty" bson:"lon,omitempty`
	Address string        `json:"address,omitempty" bson:"address,omitempty`
	Phone   string        `json:"phone,omitempty" bson:"phone,omitempty`
	Puzzles []string      `json:"puzzles,omitempty bson:"puzzles,omitempty`
}
