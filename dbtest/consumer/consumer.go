package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/sergeyyakubov/hidra2_prototype/dbtest/conf"
	"time"
	"github.com/sergeyyakubov/hidra2_prototype/dbtest/common"
)

var config conf.Config

type Pointer struct {
	ID    int `bson:"_id"`
	Value int `bson:"p"`
}

func main() {

	var configFileName string

	flag.StringVar(&configFileName, "c", "", "Config file name")
	flag.StringVar(&config.Database.Server, "s", "", "Database server")
	flag.StringVar(&config.Database.Name, "dbname", "test", "Database name")
	flag.IntVar(&config.Nthreads, "p", 1, "Number of threads to use")
	flag.IntVar(&config.Consumer.GroupID, "id", 1, "Consumer group id")
	flag.IntVar(&config.Consumer.NImagesPerRead, "n", 1, "Number of images per read")

	flag.Parse()

	if configFileName != "" {
		if err := conf.ReadConfig(configFileName, &config); err != nil {
			fmt.Println("Error reading config file " + configFileName)
			os.Exit(1)
		}
	}

	flag.Parse()

	db, err := common.ConnectDb(config)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)

	}
	defer db.Close()

	var nrecords int
	if nrecords, err = db.GetNRecords(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Running consumer with %d threads, single read: %d images, total %d records\n",
		config.Nthreads, config.Consumer.NImagesPerRead, nrecords)

	var curPointer Pointer

	db.UseCollection("groupId")

	db.DeleteAllRecords()

	err = db.IncrementField(config.Consumer.GroupID, 0, "p", &curPointer)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for i := 0; i < config.Nthreads; i++ {
		go consume(i)
	}

	lastPointer := 0
	for {
		time.Sleep(1 * time.Second)
		err = db.IncrementField(config.Consumer.GroupID, 0, "p", &curPointer)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("%d\n", curPointer.Value - lastPointer)
		lastPointer = curPointer.Value
	}
}

func consume(nthread int) {
	db, err := common.ConnectDb(config)
	defer db.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	var curPointer Pointer
	for {
		db.UseCollection("groupId")
		err := db.IncrementField(config.Consumer.GroupID, config.Consumer.NImagesPerRead, "p", &curPointer)
		if err != nil {
			fmt.Println(err)
			return
		}
		db.UseDefaultCollection()

		if config.Consumer.NImagesPerRead == 1 {
			var rec common.Record
			err = db.GetRecordByID(curPointer.Value - 1, &rec)
			if err != nil {
				fmt.Println(err)
				return
			}
			processRecord(rec)
		} else {
			var records []common.Record
			err = db.GetRecordsByIDRange(curPointer.Value - config.Consumer.NImagesPerRead, curPointer.Value - 1, &records)
			if err != nil {
				fmt.Println(err)
				return
			}
			if len(records)==0{
				return
			}
			for _,rec:=range(records){
				processRecord(rec)
			}
		}
	}
}

func processRecord(rec common.Record) {
//	fmt.Println(rec.SegmentID)
}
