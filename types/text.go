package types

import (
	"strings"

	"gopkg.in/mgo.v2/bson"
)

type Text struct {
	Content string `json:"content" bson:"content"`
	TTL     uint64 `json:"ttl,omitempty"`
}

func NewText(content string) *Text {
	return &Text{Content: content}
}

func (text *Text) Size() int {
	return len(text.Content)
}

func (text *Text) Append(content string) {
	text.Content += content
}

func (text *Text) Contains(target string) bool {
	return strings.Contains(text.Content, target)
}

func (text *Text) Equals(other *Text) bool {
	return text.Content == other.Content
}

func (text *Text) Clone() *Text {
	return &Text{Content: text.Content, TTL: text.TTL}
}

func (text *Text) Clear() {
	text.TTL = 0
	text.Content = ""
}

func (text Text) ToBSON() ([]byte, error) {
	return bson.Marshal(text)
}
