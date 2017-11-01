package jobdatabase

import (
	"time"

	"errors"

	"github.com/sergeyyakubov/dcomp/dcomp/server"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Mongodb struct {
	session *mgo.Session
	name    string
	col     string
	timeout time.Duration
	srv     *server.Server
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

func (db *Mongodb) SetServer(srv *server.Server) {
	db.srv = srv
}

func (db *Mongodb) Connect() error {
	var err error
	db.session, err = mgo.DialWithTimeout(db.srv.FullName(), db.timeout)
	return err
}

func (db *Mongodb) Close() {
	if db.session == nil {
		return
	}

	db.session.Close()
	db.session = nil
}

func (db *Mongodb) SetDefaults(name ...interface{}) {
	if len(name) > 0 {
		db.name = name[0].(string)
	}
	db.col = "jobs"
	db.timeout = 10 * time.Second
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
