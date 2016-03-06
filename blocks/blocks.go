// package blocks contains the lowest level of ipfs data structures,
// the raw block with a checksum.
package blocks

import (
	"errors"
	"fmt"

	key "github.com/ipfs/go-ipfs/blocks/key"
	mh "gx/ipfs/QmYf7ng2hG5XBtJA3tN34DQ2GUN5HNksEw1rLDkmr6vGku/go-multihash" // Block is a singular block of data in ipfs
	u "gx/ipfs/QmZNVWh8LLjAavuQ2JXuFmuYH3C11xo988vSgp7UQrTRj1/go-ipfs-util"
)

type Block struct {
	Multihash mh.Multihash
	Data      []byte
}

// NewBlock creates a Block object from opaque data. It will hash the data.
func NewBlock(data []byte) *Block {
	return &Block{Data: data, Multihash: u.Hash(data)}
}

// NewBlockWithHash creates a new block when the hash of the data
// is already known, this is used to save time in situations where
// we are able to be confident that the data is correct
func NewBlockWithHash(data []byte, h mh.Multihash) (*Block, error) {
	if u.Debug {
		chk := u.Hash(data)
		if string(chk) != string(h) {
			return nil, errors.New("Data did not match given hash!")
		}
	}
	return &Block{Data: data, Multihash: h}, nil
}

// Key returns the block's Multihash as a Key value.
func (b *Block) Key() key.Key {
	return key.Key(b.Multihash)
}

func (b *Block) String() string {
	return fmt.Sprintf("[Block %s]", b.Key())
}

func (b *Block) Loggable() map[string]interface{} {
	return map[string]interface{}{
		"block": b.Key().String(),
	}
}
