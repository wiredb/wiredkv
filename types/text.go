package types

import "gopkg.in/mgo.v2/bson"

type Text struct {
	Content string
}

func (text Text) ToBSON() ([]byte, error) {
	return bson.Marshal(text)
}
