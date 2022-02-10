package p2p

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"sync"

	//internal inports

	"github.com/pranjalpokharel7/yudhishthira/blockchain"
)

var bChain *blockchain.BlockChain
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
	Transaction []byte
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

func SendBlocks(addr string, blockchain *blockchain.BlockChain) {
	// TODO: Use ppok's json of blockchain and encode it
	// TODO: Decide on the encoding method
	var blocks = GetBlocks{
		AddrFrom: addr,
		// data:     ,
	}

	sendData(addr, blocks.data)

}

func HandleConnection(conn net.Conn, chain *blockchain.BlockChain) {
	// Reader is inderface with read method
	req, err := ioutil.ReadAll(conn)

	defer conn.Close()
	if err != nil {
		log.Panic(err)
	}

	command := req[:12]

	switch command {
	default:
		fmt.Println("Unknown command")
	}
}
