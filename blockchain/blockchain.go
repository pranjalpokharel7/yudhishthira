package blockchain

import (
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
			var genesisBlock Block

			genesisBlock.CreateGenesisBlock()
			err = ProofOfWork(&genesisBlock, DIFFICULTY)
			utility.ErrThenPanic(err)

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

	// 1) Get the hash of the last block from the chain
	err := blockchain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(LAST_HASH))
		utility.ErrThenPanic(err)

		err = item.Value(func(val []byte) error {
			lastHash = append(lastHash, val...)
			return nil
		})
		return err
	})
	utility.ErrThenPanic(err)

	// 2) Create new block with last hash pointed to the last hash key value in the database
	latestBlock.PreviousHash = lastHash
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
		fmt.Println("Blockchain iteration complete!")
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

// i.e. find unspent transaction outputs - UTXOs
func (blockchain *BlockChain) FindItemsOwned(pubKeyHash []byte) (map[string]Tx, error) {
	objectsOwned := make(map[string]Tx)
	// var objectsOwned [][]byte
	return objectsOwned, nil
}
