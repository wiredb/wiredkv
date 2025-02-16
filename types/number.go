package types

import "gopkg.in/mgo.v2/bson"

type Number struct{}

func (num Number) ToBSON() ([]byte, error) {
	return bson.Marshal(num)
}
