package database

import (
	"time"

	"errors"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Mongodb struct {
	session           *mgo.Session
	name              string
	col               string
	timeout           time.Duration
	srv               string
	ensureWriteToDisk bool
}

func (db *Mongodb) updateSession() error {
	if db.session == nil {
		return db.Connect()
	}

	if err := db.session.Ping(); err == nil {
		// nothing more to do
		return nil
	}
	db.Close()
	return db.Connect()
}

// CreateRecord changes a database record with given id. s should be is an object that
// mgo understands (go struct is OK)
func (db *Mongodb) PatchRecord(id string, s interface{}) error {
	if err := checkID(id); err != nil {
		return err
	}

	if err := db.updateSession(); err != nil {
		return err
	}

	c := db.session.DB(db.name).C(db.col)
	return c.UpdateId(bson.ObjectIdHex(id), bson.M{"$set": s})
}

// InsertRecord inserts a database record. s should be is an object that
// mgo understands (go struct is OK)
func (db *Mongodb) InsertRecord(s interface{}) error {

	if db.session == nil {
		return errors.New("database session not created")
	}

	if err := db.updateSession(); err != nil {
		return err
	}
	c := db.session.DB(db.name).C(db.col)

	return c.Insert(s)

}


// CreateRecord creates a database record with new unique id. s should be is an object that
// mgo understands (go struct is OK)
func (db *Mongodb) CreateRecord(given_id string, s interface{}) (string, error) {

	if db.session == nil {
		return "", errors.New("database session not created")
	}

	if err := db.updateSession(); err != nil {
		return "", err
	}
	c := db.session.DB(db.name).C(db.col)

	var id bson.ObjectId
	if given_id == "" {
		// create new unique id
		id = bson.NewObjectId()
	} else {
		if bson.IsObjectIdHex(given_id) {
			id = bson.ObjectIdHex(given_id)
		} else {
			return "", errors.New("Bad id format")
		}
	}

	_, err := c.UpsertId(id, s)

	if err != nil {
		return "", errors.New("Cannot add record to database: " + err.Error())
	}
	// we keep both object id for faster search and its hex representation which can be passed to clients
	// within JSON struct
	err = c.UpdateId(id, bson.M{"$set": bson.M{"_hex_id": id.Hex()}})
	if err != nil {
		return "", errors.New("Cannot update record in database: " + err.Error())
	}

	return id.Hex(), nil
}

func (db *Mongodb) SetServer(srv string) {
	db.srv = srv
}

func (db *Mongodb) Connect() error {
	var err error
	db.session, err = mgo.DialWithTimeout(db.srv, db.timeout)
	if err != nil {
		return err
	}
	db.session.SetSafe(&mgo.Safe{J:db.ensureWriteToDisk})
	return err
}

func (db *Mongodb) Close() {
	if db.session == nil {
		return
	}

	db.session.Close()
	db.session = nil
}

func (db *Mongodb) SetParams(params ...interface{}) error {
	if len(params) != 4 {
		return errors.New("Mongodb SetParams: Wrong number of params")
	}
	var ok bool
	db.name, ok = params[0].(string)
	if !ok{
		return errors.New("Mongodb: Wrong param db.name")
	}
	db.col, ok = params[1].(string)
	if !ok{
		return errors.New("Mongodb: Wrong param db.col")
	}
	db.timeout, ok = params[2].(time.Duration)
	if !ok{
		return errors.New("Mongodb: Wrong param db.timeout")
	}
	db.ensureWriteToDisk, ok = params[3].(bool)
	if !ok{
		return errors.New("Mongodb: Wrong param db.ensureWriteToDisk")
	}

	return nil
}

// GetRecords issues a request to mongodb. q should be a bson.M object or go struct with fields to match
// returns
func (db *Mongodb) GetRecords(q interface{}, res interface{}) error {

	if err := db.updateSession(); err != nil {
		return err
	}
	c := db.session.DB(db.name).C(db.col)

	return c.Find(q).All(res)

}

func (db *Mongodb) GetNRecords() (n int, err error) {

	n = 0
	if err = db.updateSession(); err != nil {
		return
	}
	c := db.session.DB(db.name).C(db.col)

	return c.Count()

}

func (db *Mongodb) DeleteAllRecords() (err error) {

	if err = db.updateSession(); err != nil {
		return
	}
	c := db.session.DB(db.name).C(db.col)

	_, err = c.RemoveAll(nil)
	return err
}



// GetAllRecords returns all records
func (db *Mongodb) GetAllRecords(res interface{}) (err error) {
	return db.GetRecords(nil, res)
}

func checkID(id string) error {
	if !bson.IsObjectIdHex(id) {
		return errors.New("wrong id")
	}
	return nil
}

func (db *Mongodb) GetRecordsByID(id string, res interface{}) error {
	if err := checkID(id); err != nil {
		return err
	}
	q := bson.M{"_id": bson.ObjectIdHex(id)}
	return db.GetRecords(q, res)
}

func (db *Mongodb) DeleteRecordByID(id string) error {

	if err := db.updateSession(); err != nil {
		return err
	}

	if err := checkID(id); err != nil {
		return err
	}
	q := bson.M{"_id": bson.ObjectIdHex(id)}

	c := db.session.DB(db.name).C(db.col)
	_, err := c.RemoveAll(q)
	return err
}
