// Package containes a database objects and functions to work with it.
// db is an interface to a specific implementation (currently implemented mongodb and mockdatabase used for tests)
package jobdatabase

import "github.com/sergeyyakubov/dcomp/dcomp/server"

type Agent interface {
	CreateRecord(string, interface{}) (string, error)
	PatchRecord(string, interface{}) error
	GetAllRecords(interface{}) error
	GetRecords(interface{}, interface{}) error
	GetRecordsByID(string, interface{}) error
	DeleteRecordByID(string) error
	Connect() error
	SetDefaults(name ...interface{})
	SetServer(*server.Server)
	Close()
}
