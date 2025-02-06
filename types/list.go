package types

type List struct{}

func (list *List) ToBSON() ([]byte, error) {
	return nil, nil
}
