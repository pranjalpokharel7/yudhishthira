package blockchain

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/dgraph-io/badger"
	"github.com/pranjalpokharel7/yudhishthira/utility"
	"github.com/pranjalpokharel7/yudhishthira/wallet"
)

// technically block chain is just a chain of blocks
// TODO: single int field to determine whether the chain is main chain or test chain
type BlockChain struct {
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
			genesisBlock := CreateGenesisBlock()
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

func (blockchain *BlockChain) AddBlock(latestBlock *Block) error {
	// check if block hash is correct
	verifiedBlockHash := latestBlock.VerifyBlockHash()
	if !verifiedBlockHash {
		return errors.New("block hash does not match")
	}

	// check if proof of work has been done on the block
	verifiedProof := latestBlock.VerifyProof()
	if !verifiedProof {
		return errors.New("proof of work hasn't been done on the block")
	}

	err := blockchain.Database.Update(func(txn *badger.Txn) error {
		latestBlockSerialized, err := latestBlock.SerializeBlockToGOB()
		utility.ErrThenPanic(err)

		err = txn.Set(latestBlock.BlockHash, latestBlockSerialized)
		utility.ErrThenPanic(err)

		err = txn.Set([]byte(LAST_HASH), latestBlock.BlockHash)
		blockchain.LastHash = latestBlock.BlockHash

		return err
	})

	return err
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

func (blockchain *BlockChain) GetLastNBlocks(n uint64) []*Block {
	var lastNBlocks []*Block

	iter := BlockChainIterator{
		CurrentHash: blockchain.LastHash,
		Database:    blockchain.Database,
	}

	for block, i := iter.GetBlockAndIter(), uint64(0); i < n && block != nil; block, i = iter.GetBlockAndIter(), i+1 {
		lastNBlocks = append(lastNBlocks, block)
	}

	return lastNBlocks
}

func (blockchain *BlockChain) GetLastNTxs(n uint64) []*Tx {
	var lastNTxs []*Tx
	var txCount uint64

	iter := BlockChainIterator{
		CurrentHash: blockchain.LastHash,
		Database:    blockchain.Database,
	}

	for block := iter.GetBlockAndIter(); txCount <= n && block != nil; block = iter.GetBlockAndIter() {
		if block.TxMerkleTree != nil {
			for _, txNode := range block.TxMerkleTree.LeafNodes {
				lastNTxs = append(lastNTxs, &txNode.Transaction)
				txCount += 1
			}
		}
	}

	return lastNTxs
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
	for block := iter.GetBlockAndIter(); block != nil && block.Height != height; block = iter.GetBlockAndIter() {
		hashes = append(hashes, block.BlockHash) // TODO: need to add in reverse order? or reverse at last?
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

func (blockchain *BlockChain) LastBlock() *Block {
	itr := &BlockChainIterator{
		CurrentHash: blockchain.LastHash,
		Database:    blockchain.Database,
	}

	block := itr.GetBlockAndIter()
	return block
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

func (blockchain *BlockChain) FindItemExists(itemHash []byte) (bool, error) {
	iter := BlockChainIterator{
		CurrentHash: blockchain.LastHash,
		Database:    blockchain.Database,
	}

	// if nil is returned then that means we reached beyond genesis block on iteration
	for block := iter.GetBlockAndIter(); block != nil; block = iter.GetBlockAndIter() {
		if block.TxMerkleTree != nil {
			for _, txNode := range block.TxMerkleTree.LeafNodes {
				if bytes.Equal(txNode.Transaction.ItemHash, itemHash) {
					// fmt.Println("Item exists in the chain beforehand")
					return true, nil
				}
			}
		}
	}

	return false, nil
}

// gets last block which contained the item
func (blockchain *BlockChain) LastBlockWithItem(itemHash []byte) (*Block, int, error) {
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
	err := fmt.Sprintf("item with hash %x does not exist", itemHash)
	return nil, -1, errors.New(err)
}

// return all transactions that contain the item
func (blockchain *BlockChain) TxsIncludingItem(itemHash []byte) []*Tx {
	var itemTxHistory []*Tx
	iter := BlockChainIterator{
		CurrentHash: blockchain.LastHash,
		Database:    blockchain.Database,
	}
	for block := iter.GetBlockAndIter(); block != nil; block = iter.GetBlockAndIter() {
		if block.TxMerkleTree != nil {
			for _, txNode := range block.TxMerkleTree.LeafNodes {
				if bytes.Equal(txNode.Transaction.ItemHash, itemHash) {
					itemTxHistory = append(itemTxHistory, &txNode.Transaction)
				}
			}
		}
	}
	return itemTxHistory
}

// get all coinbase transactions from the chain i.e. transactions in which an item was first introduced in the chain
func (blockchain *BlockChain) AllCoinBaseTxs() []*Tx {
	var tx []*Tx
	iter := BlockChainIterator{
		CurrentHash: blockchain.LastHash,
		Database:    blockchain.Database,
	}
	for block := iter.GetBlockAndIter(); block != nil; block = iter.GetBlockAndIter() {
		if block.TxMerkleTree != nil {
			for _, txNode := range block.TxMerkleTree.LeafNodes {
				if txNode.Transaction.IsCoinbase() {
					tx = append(tx, &txNode.Transaction)
				}
			}
		}
	}
	return tx
}

// get all coinbase txs by the wallet (number of items introduced into the chain)
func (blockchain *BlockChain) WalletCoinBaseTxs(walletAddress string) ([]*Tx, error) {
	var userCoinBaseTxs []*Tx
	coinBaseTxs := blockchain.AllCoinBaseTxs()
	pubKeyHash, err := wallet.PubKeyHashFromAddress(walletAddress)
	if err != nil {
		return nil, err
	}
	for _, coinbaseTx := range coinBaseTxs {
		if bytes.Equal(coinbaseTx.BuyerHash, pubKeyHash) {
			userCoinBaseTxs = append(userCoinBaseTxs, coinbaseTx)
		}
	}
	return userCoinBaseTxs, nil
}

// get available rewards for further transactions
func (blockchain *BlockChain) WalletMinedBlocks(walletAddress string) ([]*Block, error) {
	var minedBlocks []*Block
	iter := BlockChainIterator{
		CurrentHash: blockchain.LastHash,
		Database:    blockchain.Database,
	}
	pubKeyHash, err := wallet.PubKeyHashFromAddress(walletAddress)
	if err != nil {
		return nil, err
	}
	for block := iter.GetBlockAndIter(); block != nil; block = iter.GetBlockAndIter() {
		if bytes.Equal(block.Miner, pubKeyHash) {
			minedBlocks = append(minedBlocks, block)
		}
	}
	return minedBlocks, nil
}

// check if wallet has sufficient funds for coinbase transaction
func HasFundsForCoinbaseTx(walletAddress string, blockchain *BlockChain) (bool, error) {
	minedBlocks, err := blockchain.WalletMinedBlocks(walletAddress)
	if err != nil {
		return false, err
	}

	coinbaseTxsDone, err := blockchain.WalletCoinBaseTxs(walletAddress)
	if err != nil {
		return false, err
	}

	hasSufficientFunds := len(minedBlocks) > MINED_TO_SPEND_RATIO*len(coinbaseTxsDone)
	return hasSufficientFunds, nil
}

// TODO: Optimize this function
// TODO: Maybe generalize this function to get states of all items in the chain, since we're doing that anyway
func (blockchain *BlockChain) WalletOwnedItems(walletAddress string) ([]string, error) {
	var ownedItems []string
	chainItems := make(map[string]bool)

	iter := BlockChainIterator{
		CurrentHash: blockchain.LastHash,
		Database:    blockchain.Database,
	}

	pubKeyHash, err := wallet.PubKeyHashFromAddress(walletAddress)
	if err != nil {
		return nil, err
	}

	for block := iter.GetBlockAndIter(); block != nil; block = iter.GetBlockAndIter() {
		if block.TxMerkleTree != nil {
			for _, txNode := range block.TxMerkleTree.LeafNodes {
				itemHashString := hex.EncodeToString(txNode.Transaction.ItemHash)

				if _, itemRecorded := chainItems[itemHashString]; !itemRecorded {
					if bytes.Equal(txNode.Transaction.BuyerHash, pubKeyHash) {
						chainItems[itemHashString] = true
					} else if bytes.Equal(txNode.Transaction.SellerHash, pubKeyHash) {
						chainItems[itemHashString] = false
					}
				}
			}
		}
	}

	for itemHash, ownershipState := range chainItems {
		if ownershipState {
			ownedItems = append(ownedItems, itemHash)
		}
	}

	return ownedItems, nil
}
