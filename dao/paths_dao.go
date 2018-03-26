package dao

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Path struct {
	ID   bson.ObjectId `bson:"_id"`
	Path string        `bson:"path"`
	URL  string        `bson:"url"`
}
type PathsDAO struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	COLLECTION = "paths"
)

func (p *PathsDAO) Connect() {
	session, err := mgo.Dial(p.Server)
	if err != nil {
		panic(err)
	}
	db = session.DB(p.Database)
}

func (p *PathsDAO) Insert(path Path) error {
	err := db.C(COLLECTION).Insert(&path)
	return err
}

func (p *PathsDAO) FindAll() ([]Path, error) {
	var paths []Path
	err := db.C(COLLECTION).Find(bson.M{}).All(&paths)
	return paths, err
}

// Find a movie by its path
func (m *PathsDAO) FindByPath(p string) (Path, error) {
	var path Path
	err := db.C(COLLECTION).Find(bson.M{"path": p}).One(&path)
	return path, err
}
