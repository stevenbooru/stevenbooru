package dal

import (
	"encoding/gob"
	"log"

	"github.com/goincremental/dal/Godeps/_workspace/src/labix.org/v2/mgo"
	"github.com/goincremental/dal/Godeps/_workspace/src/labix.org/v2/mgo/bson"
)

func init() {
	gob.Register(ObjectID(""))
}

type iter struct {
	Iter
	iter *mgo.Iter
}

func (i *iter) Next(inter interface{}) bool {
	return i.iter.Next(inter)
}

type query struct {
	query *mgo.Query
}

func (q *query) One(i interface{}) error {
	err := q.query.One(i)
	if err == mgo.ErrNotFound {
		err = ErrNotFound
	}
	return err
}

func (q *query) All(i interface{}) error {
	err := q.query.All(i)
	if err == mgo.ErrNotFound {
		err = ErrNotFound
	}
	return err
}

func (q *query) Iter() Iter {
	i := q.query.Iter()
	return &iter{iter: i}
}

func (q *query) Sort(s ...string) Query {
	q2 := q.query.Sort(s...)
	return &query{query: q2}
}

func (q *query) Apply(change Change, result interface{}) (info *ChangeInfo, err error) {
	c := mgo.Change{
		Update:    change.Update,
		Upsert:    change.Upsert,
		Remove:    change.Remove,
		ReturnNew: change.ReturnNew,
	}
	mci, err := q.query.Apply(c, result)
	info = &ChangeInfo{}
	if mci != nil {
		info.Updated = mci.Updated
		info.Removed = mci.Removed
		info.UpsertedId = mci.UpsertedId
	}
	return
}

type collection struct {
	col *mgo.Collection
}

func (c *collection) Find(q Q) Query {
	bsonQ := c.col.Find(q)
	return &query{query: bsonQ}
}

func (c *collection) EnsureIndex(index Index) error {
	i := mgo.Index{
		Key:         index.Key,
		Background:  index.Background,
		Sparse:      index.Sparse,
		ExpireAfter: index.ExpireAfter,
	}
	return c.col.EnsureIndex(i)
}

func (c *collection) FindID(id interface{}) Query {
	q := c.col.FindId(id)
	return &query{query: q}
}

func (c *collection) RemoveID(id interface{}) error {
	return c.col.RemoveId(id)
}

func (c *collection) UpsertID(id interface{}, update interface{}) (*ChangeInfo, error) {
	return c.Upsert(bson.M{"_id": id}, update)
}

func (c *collection) SaveID(id interface{}, item interface{}) (*ChangeInfo, error) {
	return c.Upsert(bson.M{"_id": id}, bson.M{"$set": item})
}

func (c *collection) Save(selector interface{}, update interface{}) (*ChangeInfo, error) {
	return c.Upsert(selector, bson.M{"$set": update})
}

func (c *collection) Upsert(selector interface{}, update interface{}) (*ChangeInfo, error) {
	mci, err := c.col.Upsert(selector, update)
	if err != nil {
		log.Printf("Error upserting: %s\n", err)
	}
	ci := &ChangeInfo{}
	if mci != nil {
		ci.Updated = mci.Updated
		ci.Removed = mci.Removed
		ci.UpsertedId = mci.UpsertedId
	}
	log.Printf("change info %s", ci)
	return ci, err
}

func (c *collection) Insert(docs ...interface{}) error {
	return c.col.Insert(docs)
}

type database struct {
	Database
	db *mgo.Database
}

func (d *database) C(name string) Collection {
	col := d.db.C(name)
	return &collection{col: col}
}

type connection struct {
	Connection
	mgoSession *mgo.Session
}

func (c *connection) Clone() Connection {
	a := c.mgoSession.Clone()
	return &connection{mgoSession: a}
}

func (c *connection) Close() {
	c.mgoSession.Close()
}

func (c *connection) DB(name string) Database {
	db := c.mgoSession.DB(name)
	return &database{db: db}
}

type dal struct {
	DAL
}

func (d *dal) Connect(s string) (Connection, error) {
	log.Printf("Connect: %s\n", s)
	mgoSession, err := mgo.Dial(s)
	return &connection{mgoSession: mgoSession}, err
}

func NewDAL() DAL {
	return &dal{}
}

type ObjectID string

func (id ObjectID) GetBSON() (interface{}, error) {
	return bson.ObjectId(id), nil
}

func (id ObjectID) Hex() string {
	return bson.ObjectId(id).Hex()
}

func (id ObjectID) Valid() bool {
	return bson.ObjectId(id).Valid()
}

// MarshalJSON turns a dal.ObjectId into a json.Marshaller.
func (id ObjectID) MarshalJSON() ([]byte, error) {
	return bson.ObjectId(id).MarshalJSON()
}

// UnmarshalJSON turns *dal.ObjectId into a json.Unmarshaller.
func (id *ObjectID) UnmarshalJSON(data []byte) (err error) {
	a := bson.NewObjectId()
	err = a.UnmarshalJSON(data)
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}
	*id = ObjectID(a)
	return
}

func ObjectIDHex(s string) ObjectID {
	return ObjectID(bson.ObjectIdHex(s))
}

func IsObjectIDHex(s string) bool {
	return bson.IsObjectIdHex(s)
}

func NewObjectID() ObjectID {
	return ObjectID(bson.NewObjectId())
}
