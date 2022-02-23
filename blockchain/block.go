package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// the body of the block only contains the transactions
// might need a few more fields
// including the consensus number in the blockheader
type Block struct {
	Nonce        uint64 // unsigned representation for now, might allocate 64 bits later, upgrade to 64 bits if version field is removed
	Timestamp    uint64 // unix date time, string representation now, might convert to uint64 if time zones are not taken into consideration
	PreviousHash []byte // hash of previous block
	BlockHash    []byte // hash of the current block
	// transactions []transaction.Tx
}

func (blk *Block) String() string {
	var lines []string
	lines = append(lines, fmt.Sprintf("------------ Block: %x ------------", blk.BlockHash))
	lines = append(lines, fmt.Sprintf("Nonce: %d", blk.Nonce))
	lines = append(lines, fmt.Sprintf("Timestamp: %d", blk.Timestamp))
	lines = append(lines, fmt.Sprintf("Previous Hash: %x", blk.PreviousHash))
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

// cast []byte to string before marshaling
func (blk *Block) MarshalBlockToJSON() ([]byte, error) {
	block_json, err := json.Marshal(struct {
		Nonce        uint64 `json:"nonce"`
		Timestamp    uint64 `json:"timestamp"`
		PreviousHash string `json:"previous_hash"`
		BlockHash    string `json:"block_hash"`
	}{
		Nonce:        blk.Nonce,
		Timestamp:    blk.Timestamp,
		PreviousHash: hex.EncodeToString(blk.PreviousHash[:]),
		BlockHash:    hex.EncodeToString(blk.BlockHash[:]),
	})

	if err != nil {
		return nil, err
	}

	return block_json, nil
}

// this function is required because we cast []byte to string while marshaling, might remove later if affects performance
func UnmarshalJSONTOBlock(jsonData []byte) (*Block, error) {
	var unmarshalInterface map[string]interface{}
	if err := json.Unmarshal(jsonData, &unmarshalInterface); err != nil {
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
		} else {
			blk.Timestamp = uint64(v.(float64))
		}
	}

	return &blk, nil
}

// since this function does not modify the actual block properties, we remove the interface from it
func CalculateHash(blk *Block, nonce uint64) []byte {
	var buf bytes.Buffer

	blockBytes := make([]byte, 8)
	blockData := nonce ^ blk.Timestamp                   // XOR timestamp and nonce
	binary.LittleEndian.PutUint64(blockBytes, blockData) // write XORed  uint64 data to buffer
	buf.Write(blockBytes)
	buf.Write(blk.PreviousHash[:])               // write blockhash to buffer
	calculatedHash := sha256.Sum256(buf.Bytes()) // calculate hash

	return calculatedHash[:]
}

func CreateBlock() *Block {
	var blk Block
	blk.Timestamp = uint64(time.Now().Unix())
	return &blk
}

func (blk *Block) LinkPreviousHash(prevBlock *Block) {
	blk.PreviousHash = prevBlock.BlockHash
}

func (blk *Block) CreateGenesisBlock() {
	// only allow this method to be called if blockchain is empty?
	blk.Timestamp = uint64(time.Now().Unix())
	blk.PreviousHash = nil
	b_hash := sha256.Sum256([]byte(GENESIS_STRING))
	blk.BlockHash = b_hash[:]
}

func (blk *Block) AddTransactionsFromPool() {}
