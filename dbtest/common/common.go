package common

type Record struct {
	ID int `bson:"_id"`
	FName string
	BufferAddr string
	SegmentID int
	Reserv [10]string
}

