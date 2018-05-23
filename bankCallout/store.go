package bankCallout

import (
	"gopkg.in/mgo.v2"
	"fmt"
)

type Store interface {
	setCallout(chatid uint32, groupname string) error
	getCalloutGroup(chatid uint32) (string, error)
}

type mongoStore struct {
	db *mgo.Database
}

func NewMongoStore(db *mgo.Database) Store{
	return &mongoStore{
		db:db,
	}
}

type callout struct {
	ChatID uint32 `bson:"_id"`
	GroupName string
}

func (s mongoStore) setCallout(chatid uint32, groupname string) error {
	c := s.db.C("bank-callout")
	q := c.FindId(chatid)
	n, err := q.Count()
	if err != nil {
		return err
	}

	if n > 0 {
		r := callout{}
		q.One(&r)
		r.GroupName = groupname
		return c.UpdateId(chatid,r)
	}

	r := callout{
		ChatID:chatid,
		GroupName:groupname,
	}

	return c.Insert(&r)
}

func (s mongoStore) getCalloutGroup(chatid uint32) (group string, err error) {
	c := s.db.C("bank-callout")
	q := c.FindId(chatid)
	n, err := q.Count()
	if n == 0 {
		err = fmt.Errorf("no callout group has been configured for chat id %v",chatid)
		return
	}

	v := callout{}
	q.One(&v)
	group = v.GroupName
	return

}





