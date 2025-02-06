package types

type Number struct{}

func (num *Number) ToBSON() ([]byte, error) {
	return nil, nil
}
