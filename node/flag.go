package node

// Flag contains caching-related metadata about a node.
type Flag struct {
	Hash  []byte // cached hash of the node (may be nil)
	Gen   uint16 // cache generation counter
	Dirty bool   // whether the node has changes that must be written to the database
}

func NewFlag(dirty bool, gen uint16) *Flag {
	return &Flag{
		Dirty: dirty,
		Gen:   gen,
	}
}

// canUnload tells whether a node can be unloaded.
func (f *Flag) canUnload(cachegen, cachelimit uint16) bool {
	return !f.Dirty && cachegen-f.Gen >= cachelimit
}
