package repository

import (
	"github.com/go-pg/pg/orm"
)

// DataSource is an interface for go-pg DB and TX types
type DataSource interface {
	Exec(query interface{}, params ...interface{}) (orm.Result, error)
	ExecOne(query interface{}, params ...interface{}) (orm.Result, error)
	Query(model interface{}, query interface{}, params ...interface{}) (orm.Result, error)
	QueryOne(model interface{}, query interface{}, params ...interface{}) (orm.Result, error)
	Model(model ...interface{}) *orm.Query
	Select(model interface{}) error
	Insert(model ...interface{}) error
	Update(model interface{}) error
	Delete(model interface{}) error
	ForceDelete(model interface{}) error
}
