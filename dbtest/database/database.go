// Package containes a database objects and functions to work with it.
// db is an interface to a specific implementation (currently implemented mongodb and mockdatabase used for tests)
package database

type Agent interface {
	InsertRecord(s interface{}) error
	CreateRecord(string, interface{}) (string, error)
	PatchRecord(string, interface{}) error
	GetAllRecords(interface{}) error
	GetNRecords() (int,error)
	GetRecords(interface{}, interface{}) error
	GetRecordsByID(string, interface{}) error
	DeleteRecordByID(string) error
	DeleteAllRecords() error
	Connect() error
	SetParams(params ...interface{}) error
	SetServer(string)
	Close()
}
