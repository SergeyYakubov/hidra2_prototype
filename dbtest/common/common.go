package common

import (
	"github.com/sergeyyakubov/hidra2_prototype/dbtest/database"
	"github.com/sergeyyakubov/hidra2_prototype/dbtest/conf"
	"time"
)

type Record struct {
	ID int `bson:"_id"`
	FName string
	BufferAddr string
	SegmentID int
	Reserv [10]string
}

func ConnectDb(config conf.Config) (db database.Agent, err error) {

	db = new(database.Mongodb)
	db.SetServer(config.Database.Server)
	db.SetParams(config.Database.Name, "records",10*time.Second,config.Database.EnsureDiskWrite)

	if err = db.Connect(); err != nil {
		return nil, err
	}
	return db, nil
}
