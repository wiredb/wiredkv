package types

import "gopkg.in/mgo.v2/bson"

type Text struct {
	Content string `json:"content" bson:"content"`
}

func (text Text) ToBSON() ([]byte, error) {
	return bson.Marshal(text)
}
