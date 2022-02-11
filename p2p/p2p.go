package p2p

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"sync"

	//internal inports

	"github.com/pranjalpokharel7/yudhishthira/blockchain"
	"github.com/pranjalpokharel7/yudhishthira/transaction"
)

// TODO: test everything

// var bChain *blockchain.BlockChain
var mutex sync.Mutex

const (
	UNNAMED              = 0x0    // not a full node, It may not be able to provide any data except for the transactions it originates.
	NODE_NETWORK         = 0x01   // full node, can be asked for full blocks
	NODE_NETWORK_LIMITED = 0x0400 // same as node network but node has at least last (some amount of blocks to be decided)
	commandLength        = 12     // command will have 12 bytes
	protocol             = "tcp"
)

var (
	knownNodes  []string
	nodeAddress string
)

// const for types
// using integer rather than strings
// well may need to serialize this too
const (
	BLOCK_TYPE   = 0
	TX_TYPE      = 1
	VERSION_TYPE = 2
	INV_TYPE     = 3
)

// wrapper struct to send a block
type Block struct {
	AddrFrom string
	Block    []byte
}

// a version message that contains the version of the chain in the node
type Version struct {
	Timestamp   uint64
	AddressFrom string
	Height      int32
}

// contains all the address of the connected nodes
type Address struct {
	AddrList []string
}

// request a particular data object from another node
// Response to get data can be a tx, block,
//TODO: add other if required
type GetData struct {
	AddrFrom string
	Type     int32
	data     []byte // needs to be a serialized array of bytes
}

// provides the block header hashes from a particular point
type GetBlocks struct {
	AddrFrom string
	data     []byte
}

// transaction wrapper
type Tx struct {
	AddrFrom    string
	Transaction []byte //
}

// For details follow this link: https://developer.bitcoin.org/reference/p2p_networking.html#inv
type Inv struct {
	AddrFrom string
	Type     int32    // specify what type of inventory are we sending
	data     [][]byte // 2D array of byte, each byte array contains a transaction
}

func CommandToBytes(cmd string) []byte {
	var bytes []byte

	for i, c := range cmd {
		bytes[i] = byte(c)
	}

	return bytes[:]
}

func BytesToCommand(bytes []byte) string {
	var cmd []byte

	for _, b := range bytes {
		if b != 0x0 {
			cmd = append(cmd, b)
		}
	}

	return fmt.Sprintf("%s", cmd)
}

// function to send all types of serialized data
// will be called from other functin for each specialized function
func sendData(addr string, data []byte) {
	conn, err := net.Dial(protocol, addr)

	if err != nil {
		fmt.Printf("Node %s is not available\n", addr)
		var updatedNodes []string
		// if the address is not available, remove that node
		for _, node := range knownNodes {
			if node != addr {
				updatedNodes = append(updatedNodes, node)
			}
		}

		knownNodes = updatedNodes

		return
	}

	defer conn.Close()

	// send the data to the connection
	_, err = io.Copy(conn, bytes.NewReader(data))

	if err != nil {
		log.Panic(err)
	}

}

// Sends get block request to another node
func SendGetBlocks(addr string) {
	var blocks = GetBlocks{
		AddrFrom: nodeAddress,
	}
	// first 12 character is command, rest is the payload
	// check this link for more details
	//READ: https://developer.bitcoin.org/reference/p2p_networking.html#headers
	info := append(CommandToBytes("getblocks"), GobEncode(blocks)...)

	sendData(addr, info)

}

// sends all the blocks
func SendBlock(addr string, block *blockchain.Block) {
	data, err := block.MarshalBlockToJSON()
	if err != nil {
		fmt.Printf("Block serialization error: %s\n", err)
		return
	}
	var blocks = Block{
		AddrFrom: nodeAddress,
		Block:    data,
	}
	info := append(CommandToBytes("block"), GobEncode(blocks)...)

	sendData(addr, info)
}

// sends all the known address
func SendAddress(addr string, block *blockchain.Block) {
	address := Address{AddrList: knownNodes}

	info := append(CommandToBytes("address"), GobEncode(address)...)

	sendData(addr, info)
}

//sends get data request, data can be of any type
// here type is represented by id
func setGetData(addr string, kind int32, id []byte) {
	data := GobEncode(GetData{
		AddrFrom: nodeAddress,
		Type:     int32(kind),
		data:     id,
	})

	data = append(CommandToBytes("getdata"), data...)
	sendData(addr, data)
}

// send a particular transaction to the given address
func setTx(addr string, tx transaction.Tx) {
	data := GobEncode(Tx{
		AddrFrom:    nodeAddress,
		Transaction: tx.Serialize(),
	})

	data = append(CommandToBytes("getdata"), data...)
	sendData(addr, data)
}

func SendVersion(addr string, bChain *blockchain.BlockChain) {
	height := bChain.GetHeight()
	data := GobEncode(Version{
		AddressFrom: nodeAddress,
		Height:      height,
	})

	data = append(CommandToBytes("getversion"), data...)

	sendData(addr, data)
}

//transmits one or more inventories of objects known to the transmitting peer.
// The receiving peer can compare the inventories from an “inv” message against the inventories it has already seen, and then use a follow-up message to request unseen objects.
// For more info: https://developer.bitcoin.org/reference/p2p_networking.html#inv
func sendInv(addr string, kind int32, inventories [][]byte) {
	data := GobEncode(Inv{
		AddrFrom: nodeAddress,
		Type:     int32(kind),
		data:     inventories,
	})

	data = append(CommandToBytes("inv"), data...)
	sendData(addr, data)
}

func HandleConnection(conn net.Conn, chain *blockchain.BlockChain) {
	// Reader is inderface with read method
	req, err := ioutil.ReadAll(conn)

	defer conn.Close()
	if err != nil {
		log.Panic(err)
	}

	// get all the required commands
	command := req[:12]

	switch command {
	default:
		fmt.Println("Unknown command")
	}
}

func GobEncode(data interface{}) []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
