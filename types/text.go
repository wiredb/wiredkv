package types

type Text struct{}

func (text *Text) ToBSON() []byte {
	return nil
}
