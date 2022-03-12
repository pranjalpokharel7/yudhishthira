package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/pranjalpokharel7/yudhishthira/wallet"
)

// the body of the block only contains the transactions
// might need a few more fields
// including the consensus number in the blockheader
type Block struct {
	Nonce        uint64      `json:"nonce"`         // unsigned representation for now, might allocate 64 bits later, upgrade to 64 bits if version field is removed
	Height       uint64      `json:"height"`        // current block height
	Timestamp    uint64      `json:"timestamp"`     // unix date time, string representation now, might convert to uint64 if time zones are not taken into consideration
	Difficulty   uint64      `json:"difficulty"`    // difficulty based on tx sum
	BlockHash    HexByte     `json:"block_hash"`    // hash of the current block
	PreviousHash HexByte     `json:"previous_hash"` // hash of previous block
	Miner        HexByte     `json:"miner"`         // address of block miner
	TxMerkleTree *MerkleTree `json:"merkle_tree"`   // merkel tree for transactions
}

func (blk *Block) String() string {
	var lines []string
	lines = append(lines, fmt.Sprintf("------------ Block: %x ------------", blk.BlockHash))
	lines = append(lines, fmt.Sprintf("Nonce: %d", blk.Nonce))
	lines = append(lines, fmt.Sprintf("Height: %d", blk.Height))
	lines = append(lines, fmt.Sprintf("Miner: %x", blk.Miner))
	lines = append(lines, fmt.Sprintf("Timestamp: %d", blk.Timestamp))
	lines = append(lines, fmt.Sprintf("Previous Hash: %x", blk.PreviousHash))

	// just print the merkle json for now, no point in worrying too much about pretty printing
	merkleJSON, _ := json.MarshalIndent(blk.TxMerkleTree, "", "\t")
	lines = append(lines, fmt.Sprintf("Merkle Tree: %s", merkleJSON))
	return strings.Join(lines, "\n")
}

func (blk *Block) SerializeBlockToGOB() ([]byte, error) {
	var encoded bytes.Buffer
	err := gob.NewEncoder(&encoded).Encode(blk)
	return encoded.Bytes(), err
}

func DeserializeBlockFromGOB(serializedBlock []byte) (*Block, error) {
	var blk Block
	err := gob.NewDecoder(bytes.NewReader(serializedBlock)).Decode(&blk)
	return &blk, err
}

// this function is required because we cast []byte to string while marshaling, might remove later if affects performance
func UnmarshalBlockFromJSON(jsonData []byte) (*Block, error) {
	var unmarshalInterface map[string]interface{}
	err := json.Unmarshal(jsonData, &unmarshalInterface)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var blk Block
	for k, v := range unmarshalInterface {
		if k == "nonce" {
			blk.Nonce = uint64(v.(float64))
		} else if k == "block_hash" {
			blk.BlockHash = []byte(v.(string))
		} else if k == "previous_hash" {
			blk.PreviousHash = []byte(v.(string))
		} else if k == "height" {
			blk.Height = uint64(v.(float64))
		} else if k == "timestamp" {
			blk.Timestamp = uint64(v.(float64))
		} else {
			blk.TxMerkleTree, err = UnmarshalMerkleFromInterface(v.(map[string]interface{}))
		}
	}

	return &blk, err
}

func CreateBlock() *Block {
	var blk Block
	blk.Timestamp = uint64(time.Now().Unix())
	return &blk
}

func CreateGenesisBlock() *Block {
	blk_hash := sha256.Sum256([]byte(GENESIS_STRING))
	// leave most fields empty for now
	blk := Block{
		Timestamp: GENESIS_TIMESTAMP,
		Height:    0,
		BlockHash: blk_hash[:],
	}
	return &blk
}

func CalculateHashEmptyBlock(blk *Block, nonce uint64) []byte {
	var buf bytes.Buffer

	blockBytes := make([]byte, 8)
	blockData := nonce ^ blk.Timestamp                   // XOR timestamp and nonce
	binary.LittleEndian.PutUint64(blockBytes, blockData) // write XORed  uint64 data to buffer
	buf.Write(blockBytes)
	buf.Write(blk.PreviousHash[:]) // write blockhash to buffer

	calculatedHash := sha256.Sum256(buf.Bytes()) // calculate hash

	return calculatedHash[:]
}

// since this function does not modify the actual block properties, we remove the interface from it
// TODO: gob encode and hash using only required fields, as done for transaction
func CalculateHashNonEmptyBlock(blk *Block, nonce uint64) []byte {
	var buf bytes.Buffer

	blockBytes := make([]byte, 8)
	blockData := nonce ^ blk.Timestamp                   // XOR timestamp and nonce
	binary.LittleEndian.PutUint64(blockBytes, blockData) // write XORed  uint64 data to buffer
	buf.Write(blockBytes)
	buf.Write(blk.PreviousHash[:])             // write blockhash to buffer
	buf.Write(blk.TxMerkleTree.Root.HashValue) // write merkel root hash to buffer

	calculatedHash := sha256.Sum256(buf.Bytes()) // calculate hash

	return calculatedHash[:]
}

func (blk *Block) MineBlock(chain *BlockChain, wlt *wallet.Wallet) error {
	var lastHash []byte
	var lastBlock *Block

	err := chain.Database.View(func(txn *badger.Txn) error {
		lastHashQuery, err := txn.Get([]byte(LAST_HASH))
		if err != nil {
			return err
		}

		err = lastHashQuery.Value(func(val []byte) error {
			lastHash = append(lastHash, val...)
			return nil
		})

		lastBlockQuery, err := txn.Get(lastHash)
		if err != nil {
			return err
		}

		err = lastBlockQuery.Value(func(val []byte) error {
			lastBlock, err = DeserializeBlockFromGOB(val)
			return err
		})
		return err
	})

	if err != nil {
		return err
	}

	blk.PreviousHash = lastHash
	blk.Height = lastBlock.Height + 1

	// create function to calculate difficulty later based on txsum?
	blk.Difficulty = blk.Height%2016 + 1 // block difficulty changes every 2016 blocks, just like bitcoin
	ProofOfWork(blk)

	// add miner address after proof of work is done
	minerAddress, err := wallet.PubKeyHashFromAddress(string(wlt.Address))
	if err != nil {
		return err
	}
	blk.Miner = minerAddress

	return nil
}

// add transactions from pool to block as merkle tree
func (blk *Block) AddTransactionsToBlock(txPool []Tx) error {
	var tree *MerkleTree
	tree, err := CreateMerkleTree(txPool, tree)
	if err != nil {
		return err
	}
	blk.TxMerkleTree = tree
	return nil
}

func (block *Block) TxSum() uint64 {
	var txSum uint64
	for _, txNode := range block.TxMerkleTree.LeafNodes {
		txSum += txNode.Transaction.Amount
	}
	return txSum
}

func (blk *Block) VerifyBlockHash() bool {
	var calculatedHash []byte
	if blk.IsEmpty() {
		calculatedHash = CalculateHashEmptyBlock(blk, blk.Nonce)
	} else {
		calculatedHash = CalculateHashNonEmptyBlock(blk, blk.Nonce)
	}
	return bytes.Equal(calculatedHash, blk.BlockHash)
}

func (blk *Block) IsEmpty() bool {
	return blk.TxMerkleTree == nil
}
