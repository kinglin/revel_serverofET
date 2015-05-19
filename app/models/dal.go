package models

import (
	"github.com/robfig/revel"
	"gopkg.in/mgo.v2"
)

type Dal struct {
	session *mgo.Session
}

func NewDal() (*Dal, error) {
	revel.Config.SetSection("db")
	ip, found := revel.Config.String("ip")
	if !found {
		revel.ERROR.Fatal("Cannot load database ip from app.conf")
	}

	session, err := mgo.Dial(ip)
	if err != nil {
		return nil, err
	}

	return &Dal{session}, nil
}

func (d *Dal) Close() {
	d.session.Close()
}
