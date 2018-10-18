package node

type Node interface {
	Find([]byte) ([]byte, Node, bool)
	Put([]byte, Node) (Node, bool)
	Delete([]byte) (Node, bool)
	Cache() ([]byte, bool)
	CanUpload(gen, limit uint16) bool
}
