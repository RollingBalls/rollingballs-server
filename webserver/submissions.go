package webserver

type PositionSubmission struct {
	Lat float32 `binding:"required"`
	Lon float32 `binding:"required"`
}

type GuessingSubmission struct {
	Id int `binding:"required"`
}
