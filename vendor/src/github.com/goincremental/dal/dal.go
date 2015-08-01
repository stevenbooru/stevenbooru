package dal

import (
	"errors"
	"time"
)

var ErrNotFound = errors.New("ErrNotFound")

type Q map[string]interface{}

type DAL interface {
	Connect(string) (Connection, error)
	IsObjectIdHex(string) bool
}

type Connection interface {
	Clone() Connection
	Close()
	DB(s string) Database
}

type Database interface {
	C(string) Collection
}

type Collection interface {
	Find(Q) Query
	EnsureIndex(Index) error
	FindID(interface{}) Query
	RemoveID(interface{}) error
	UpsertID(interface{}, interface{}) (*ChangeInfo, error)
	Upsert(interface{}, interface{}) (*ChangeInfo, error)
	Insert(...interface{}) error
	Save(interface{}, interface{}) (*ChangeInfo, error)
	SaveID(interface{}, interface{}) (*ChangeInfo, error)
}

type Query interface {
	One(interface{}) error
	All(interface{}) error
	Sort(...string) Query
	Iter() Iter
	Apply(Change, interface{}) (*ChangeInfo, error)
}

type Iter interface {
	Next(interface{}) bool
}

type Index struct {
	Key         []string
	Background  bool
	Sparse      bool
	ExpireAfter time.Duration
}

type ChangeInfo struct {
	Updated    int
	Removed    int
	UpsertedId interface{}
}

type Change struct {
	Update    interface{}
	Upsert    bool
	Remove    bool
	ReturnNew bool
}
