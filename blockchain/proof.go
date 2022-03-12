package blockchain

import (
	"encoding/hex"
	"errors"
	"strings"
)

const DIFFICULTY = 3 // make difficulty dynamic later on

func containsLeadingZeroes(hash []byte, difficulty uint64) bool {
	var hexRepresentation string = hex.EncodeToString(hash[:])
	var leadingZeroes string = strings.Repeat("0", int(difficulty))
	return hexRepresentation[0:difficulty] == leadingZeroes
}

func ProofOfWork(blk *Block) error {
	HashStrategy := CalculateHashNonEmptyBlock
	if blk.IsEmpty() {
		HashStrategy = CalculateHashEmptyBlock
	}
	for i := uint64(0); i < MAX_ITERATIONS_POW; i++ { // arbitrary 1000 to prevent potential endless loop
		hash := HashStrategy(blk, i)
		if containsLeadingZeroes(hash, blk.Difficulty) {
			blk.BlockHash = hash
			blk.Nonce = i
			return nil
		}
	}
	return errors.New("proof of work could not be calculated within the given number of iterations")
}

func (blk *Block) VerifyProof() bool {
	return containsLeadingZeroes(blk.BlockHash, blk.Difficulty)
}
