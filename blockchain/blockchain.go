package blockchain

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/dgraph-io/badger"
	"github.com/pranjalpokharel7/yudhishthira/utility"
)

// technically block chain is just a chain of blocks
// TODO: single int field to determine whether the chain is main chain or test chain
type BlockChain struct {
	// Difficulty uint32 // --. moved to proof since it seems relevant there
	Database *badger.DB
	LastHash []byte
}

type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func InitBlockChain() *BlockChain {
	var lastHash []byte
	db, err := badger.Open(badger.DefaultOptions(DB_PATH))
	utility.ErrThenPanic(err)

	// to perform read-write operations, use Update
	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte(LAST_HASH)); err == badger.ErrKeyNotFound {
			// no blocks in the blockchain yet, need to add genesis block
			// TODO: separate this into a different function, since we need to run this just once in production

			genesisBlock := CreateGenesisBlock()
			// err = ProofOfWork(genesisBlock, DIFFICULTY)
			// utility.ErrThenPanic(err)

			genesisSerialized, err := genesisBlock.SerializeBlockToGOB()
			utility.ErrThenPanic(err)

			err = txn.Set(genesisBlock.BlockHash, genesisSerialized)
			utility.ErrThenPanic(err)

			err = txn.Set([]byte(LAST_HASH), genesisBlock.BlockHash)

			lastHash = append(lastHash, genesisBlock.BlockHash...)
			return err
		}

		// run a get transaction to get the last hash of the chain
		item, err := txn.Get([]byte(LAST_HASH))
		utility.ErrThenPanic(err)

		// access value from key-value pair
		err = item.Value(func(val []byte) error {
			lastHash = append(lastHash, val...)
			return nil
		})
		return err
	})

	utility.ErrThenPanic(err)

	blockchain := &BlockChain{db, lastHash}
	return blockchain
}

func (blockchain *BlockChain) AddBlock(latestBlock *Block) {
	var lastHash []byte
	var lastBlock *Block

	// 1) Get the hash of the last block from the chain
	err := blockchain.Database.View(func(txn *badger.Txn) error {
		lastHashQuery, err := txn.Get([]byte(LAST_HASH))
		utility.ErrThenPanic(err)

		err = lastHashQuery.Value(func(val []byte) error {
			lastHash = append(lastHash, val...)
			return nil
		})

		lastBlockQuery, err := txn.Get(lastHash)
		utility.ErrThenPanic(err)

		err = lastBlockQuery.Value(func(val []byte) error {
			lastBlock, err = DeserializeBlockFromGOB(val)
			return err
		})
		return err
	})
	utility.ErrThenPanic(err)

	// 2) Create new block with last hash pointed to the last hash key value in the database
	latestBlock.PreviousHash = lastHash
	latestBlock.Height = lastBlock.Height + 1
	ProofOfWork(latestBlock, DIFFICULTY) // TODO: create an abstraction methodf MineBlock(), POW can only be run after linking previous hash

	err = blockchain.Database.Update(func(txn *badger.Txn) error {
		latestBlockSerialized, err := latestBlock.SerializeBlockToGOB()
		utility.ErrThenPanic(err)

		err = txn.Set(latestBlock.BlockHash, latestBlockSerialized)
		utility.ErrThenPanic(err)

		err = txn.Set([]byte(LAST_HASH), latestBlock.BlockHash)
		blockchain.LastHash = latestBlock.BlockHash

		return err
	})

	utility.ErrThenPanic(err)
}

// return the last block from the chain and iterator backwards in the chain
func (iter *BlockChainIterator) GetBlockAndIter() *Block {
	if iter.CurrentHash == nil {
		// fmt.Println("Blockchain iteration complete!")
		return nil
	}
	var block *Block

	// to perform read only transaction, use the View method
	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		utility.ErrThenPanic(err)
		err = item.Value(func(val []byte) error {
			block, err = DeserializeBlockFromGOB(val)
			return err
		})
		return err
	})

	utility.ErrThenPanic(err)
	iter.CurrentHash = block.PreviousHash
	return block
}

func (chain *BlockChain) GetChainHeight() (uint64, error) {
	var block *Block

	if chain.Database == nil {
		return 0, nil
	}

	// to perform read only transaction, use the View method
	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(chain.LastHash)
		utility.ErrThenPanic(err)
		err = item.Value(func(val []byte) error {
			block, err = DeserializeBlockFromGOB(val)
			return err
		})
		return err
	})

	return block.Height, err
}

func (blockchain *BlockChain) GetHeight() uint64 {
	height, _ := blockchain.GetChainHeight()

	return height
}

func (blockchain *BlockChain) GetBlockHashes(blockHash []byte) [][]byte {
	var hashes [][]byte
	var hashesInOrder [][]byte

	iter := BlockChainIterator{
		CurrentHash: blockchain.LastHash,
		Database:    blockchain.Database,
	}

	// we only need heights after a certain block and not the block with the matching itself
	block := iter.GetBlockAndIter()
	for block != nil && !bytes.Equal(block.BlockHash, blockHash) {
		hashes = append(hashes, block.BlockHash)
		block = iter.GetBlockAndIter()
	}

	for i := len(hashes) - 1; i >= 0; i-- {
		hashesInOrder = append(hashesInOrder, hashes[i])
	}

	return hashesInOrder
}

func (blockchain *BlockChain) GetBlockHashesFromHeight(height uint64) [][]byte {
	var hashes [][]byte
	var hashesInOrder [][]byte

	iter := BlockChainIterator{
		CurrentHash: blockchain.LastHash,
		Database:    blockchain.Database,
	}

	// we only need heights after a certain block and not the block with the matching itself
	block := iter.GetBlockAndIter()
	for block.Height != height {
		hashes = append(hashes, block.BlockHash) // TODO: need to add in reverse order? or reverse at last?
		block = iter.GetBlockAndIter()
	}

	for i := len(hashes) - 1; i >= 0; i-- {
		hashesInOrder = append(hashesInOrder, hashes[i])
	}

	return hashesInOrder
}

//return aa block with a particular hash
func (blockchain *BlockChain) GetBlock(blockhash []byte) (*Block, error) {
	itr := &BlockChainIterator{
		CurrentHash: blockchain.LastHash,
		Database:    blockchain.Database,
	}

	for b := itr.GetBlockAndIter(); b != nil; b = itr.GetBlockAndIter() {
		if bytes.Equal(blockhash, b.BlockHash) {
			return b, nil
		}
	}

	err := errors.New("Block not found")
	return nil, err
}

func (blockchain *BlockChain) PrintChain() {
	iter := BlockChainIterator{
		CurrentHash: blockchain.LastHash,
		Database:    blockchain.Database,
	}
	block := iter.GetBlockAndIter()
	for block != nil {
		fmt.Println(block)
		block = iter.GetBlockAndIter()
	}
}

// i.e. find unspent transaction outputs - UTXOs
func (blockchain *BlockChain) FindItemsOwned(pubKeyHash []byte) (map[string]Tx, error) {
	objectsOwned := make(map[string]Tx)
	// var objectsOwned [][]byte
	return objectsOwned, nil
}

func (blockchain *BlockChain) FindItemExists(itemHash []byte) (bool, error) {
	iter := BlockChainIterator{
		CurrentHash: blockchain.LastHash,
		Database:    blockchain.Database,
	}

	// if nil is returned then that means we reached the genesis block on iteration
	for block := iter.GetBlockAndIter(); block != nil; block = iter.GetBlockAndIter() {
		for _, txNode := range block.TxMerkleTree.LeafNodes {
			if bytes.Equal(txNode.Transaction.ItemHash, itemHash) {
				// fmt.Println("Item exists in the chain beforehand")
				return true, nil
			}
		}
	}

	return false, nil
}

func (blockchain *BlockChain) GetLastBlockWithItem(itemHash []byte) (*Block, int, error) {
	iter := BlockChainIterator{
		CurrentHash: blockchain.LastHash,
		Database:    blockchain.Database,
	}
	for block := iter.GetBlockAndIter(); block != nil; block = iter.GetBlockAndIter() {
		if block.TxMerkleTree != nil {
			for txIndex, txNode := range block.TxMerkleTree.LeafNodes {
				if bytes.Equal(txNode.Transaction.ItemHash, itemHash) {
					return block, txIndex, nil
				}
			}
		}
	}
	return nil, -1, errors.New("block with item does not exist")
}

// return blocks that contains the item
func GetAllBlocksWithItem() {

}
