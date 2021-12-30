package blockchain

// technically block chain is just a chain of blocks
// a single int field to determine whether the chain is main chain or test chain
type BlockChain struct {
	chainType CHAIN_TYPE
	Blocks    []Block
}

// only has the main chain for now
func createGenesisBlock() {}

// load block chain from JSON format
func loadBlockChain() {}

// broadcast your changes to the chain from this function
func updateBlockChain() {}
