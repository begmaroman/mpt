package mptv2

type BatchParameter struct {
	Type  string
	Key   []byte
	Value []byte
}

type Database interface {
	Get(key []byte) ([]byte, error)
	Put(key, value []byte) ([]byte, error)
	Del(key []byte) error
	Batch(params []*BatchParameter) error
}
