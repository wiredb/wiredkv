package types

import "gopkg.in/mgo.v2/bson"

// curl -X PUT http://192.168.31.221:2668/tables/key-03 \
//      -H "Content-Type: application/json" \
//      -H "Auth: 11111" \
//      -d '{
//        "table": {
//          "is_valid": false,
//          "items": [
//            {"id": 1, "name": "Item 1"},
//            {"id": 2, "name": "Item 2"}
//          ],
//          "meta": {
//            "version": "2.0",
//            "author": "John Doe"
//          }
//        },
//        "TTL": 10,
//      }'

// {"code":200,"message":"request processed successfully!"}

type Tables struct {
	Table map[string]any `json:"table"`
	TTL   uint64         `json:"ttl,omitempty"`
}

func (tab Tables) ToBSON() ([]byte, error) {
	return bson.Marshal(tab)
}
