package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/sergeyyakubov/hidra2_prototype/dbtest/conf"
	"github.com/sergeyyakubov/hidra2_prototype/dbtest/database"
	"time"
	"github.com/sergeyyakubov/hidra2_prototype/dbtest/common"
)

var config conf.Config

func main() {

	var configFileName string

	flag.StringVar(&configFileName, "c", "", "Config file name")
	flag.IntVar(&config.Nthreads, "p", 1, "Number of threads to use")
	flag.StringVar(&config.Database.Server, "s", "", "Database server")
	flag.StringVar(&config.Database.Name, "dbname", "testproducer", "Database name")
	flag.BoolVar(&config.Database.EnsureDiskWrite, "j", false, "Ensure write to disk")


	flag.Parse()

	if configFileName != "" {
		if err := conf.ReadConfig(configFileName, &config); err != nil {
			fmt.Println("Error reading config file " + configFileName)
			os.Exit(1)
		}
	}

	flag.Parse()

	db, err := connectDb()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)

	}
	defer db.Close()

	if err := db.DeleteAllRecords(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Running with %d threads, ensure write to disk: %t\n",config.Nthreads,config.Database.EnsureDiskWrite )

	for i := 0; i < config.Nthreads; i++ {
		go produce(i, config.Nthreads)
	}

	nlast := 0
	for {
		time.Sleep(1 * time.Second)
		n, err := db.GetNRecords()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)

		}
		fmt.Printf("%d\n", n - nlast)
		nlast = n
	}
}

func connectDb() (db database.Agent, err error) {

	db = new(database.Mongodb)

	db.SetServer(config.Database.Server)
	db.SetParams(config.Database.Name, "records",10*time.Second,config.Database.EnsureDiskWrite)

	if err = db.Connect(); err != nil {
		return nil, err
	}
	return db, nil
}

func produce(start, shift int) {
	curID := start
	db, err := connectDb()
	defer db.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		var rec common.Record
		rec.ID = curID
		rec.BufferAddr = "126.567.344.346:45600"
		rec.SegmentID = 12345
		rec.FName = "/data/tztf/sdfsdf/sdfsdf"
		rec.Reserv[0] = "fdgfdgdfgsdfgsdgsdkfjgbdsibgiub"

		err := db.InsertRecord(rec)
		if err != nil {
			fmt.Println(err)
			return
		}
		if curID > 150000 {
			break
		}
		curID += shift
	}

}