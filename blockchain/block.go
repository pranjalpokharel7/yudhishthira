package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/pranjalpokharel7/yudhishthira/utility"
)

// the body of the block only contains the transactions
// might need a few more fields
// including the consensus number in the blockheader
type Block struct {
	Nonce        uint64      `json:"nonce"`         // unsigned representation for now, might allocate 64 bits later, upgrade to 64 bits if version field is removed
	Height       uint64      `json:"height"`        // current block height
	Timestamp    uint64      `json:"timestamp"`     // unix date time, string representation now, might convert to uint64 if time zones are not taken into consideration
	BlockHash    HexByte     `json:"block_hash"`    // hash of the current block
	PreviousHash HexByte     `json:"previous_hash"` // hash of previous block
	TxMerkleTree *MerkleTree `json:"merkle_tree"`
}

func (blk *Block) String() string {
	var lines []string
	lines = append(lines, fmt.Sprintf("------------ Block: %x ------------", blk.BlockHash))
	lines = append(lines, fmt.Sprintf("Nonce: %d", blk.Nonce))
	lines = append(lines, fmt.Sprintf("Height: %d", blk.Height))
	lines = append(lines, fmt.Sprintf("Timestamp: %d", blk.Timestamp))
	lines = append(lines, fmt.Sprintf("Previous Hash: %x", blk.PreviousHash))

	// just print the merkle json for now, no point in worrying too much about pretty printing
	merkleJSON, _ := json.Marshal(blk.TxMerkleTree)
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

// since this function does not modify the actual block properties, we remove the interface from it
func CalculateHash(blk *Block, nonce uint64) []byte {
	var buf bytes.Buffer

	blockBytes := make([]byte, 8)
	blockData := nonce ^ blk.Timestamp                   // XOR timestamp and nonce
	binary.LittleEndian.PutUint64(blockBytes, blockData) // write XORed  uint64 data to buffer
	buf.Write(blockBytes)
	buf.Write(blk.PreviousHash[:]) // write blockhash to buffer
	if blk.TxMerkleTree == nil {
		utility.ErrThenPanic(errors.New("no transactions added to the block yet"))
	}
	buf.Write(blk.TxMerkleTree.Root.HashValue)   // write merkel root hash to buffer
	calculatedHash := sha256.Sum256(buf.Bytes()) // calculate hash

	return calculatedHash[:]
}

func CreateBlock() *Block {
	var blk Block
	blk.Timestamp = uint64(time.Now().Unix())
	return &blk
}

func CreateGenesisBlock() *Block {
	var blk Block
	blk.Timestamp = uint64(time.Now().Unix())
	blk.PreviousHash = nil
	blk.Height = 0
	blk.TxMerkleTree = nil
	b_hash := sha256.Sum256([]byte(GENESIS_STRING))
	blk.BlockHash = b_hash[:]
	return &blk
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
