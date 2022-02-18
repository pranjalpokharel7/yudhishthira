package p2p

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
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
// TODO: Implement Initial Block Download

var mutex sync.Mutex

const (
	UNNAMED              = 0x0    // not a full node, It may not be able to provide any data except for the transactions it originates.
	NODE_NETWORK         = 0x01   // full node, can be asked for full blocks
	NODE_NETWORK_LIMITED = 0x0400 // same as node network but node has at least last (some amount of blocks to be decided)
	commandLength        = 12     // command will have 12 bytes
	protocol             = "tcp"
)

var (
	knownNodes  []string // list of all the knownNodes
	nodeAddress string   // address of this node

	// here string is the transaction id and it point to the actual transaction
	memoryPool = make(map[string]transaction.Tx)
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

type MESSAGE_TYPE int

// wrapper struct to send a block
type Block struct {
	AddrFrom string
	Block    []byte
}

// a version message that contains the version of the chain in the node
type Version struct {
	Timestamp   uint64
	AddressFrom string
	Height      int
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
	Type     MESSAGE_TYPE
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

// For details follow this link: https://developer.bitcoin.org/reference/p2p_networking.html#inv
type Inv struct {
	AddrFrom string
	Type     MESSAGE_TYPE // specify what type of inventory are we sending
	data     [][]byte     // 2D array of byte, each byte array contains a transaction
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
// will be called from other function for each specialized function
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

func SendAddress(addr string, block *blockchain.Block) {
	address := Address{AddrList: knownNodes}

	info := append(CommandToBytes("address"), GobEncode(address)...)

	sendData(addr, info)
}

//sends get data request, data can be of any type
// here type is represented by id
func sendGetData(addr string, kind MESSAGE_TYPE, id []byte) {
	data := GobEncode(GetData{
		AddrFrom: nodeAddress,
		Type:     kind,
		data:     id,
	})

	data = append(CommandToBytes("getdata"), data...)
	sendData(addr, data)
}

// send a particular transaction to the given address
func sendTx(addr string, tx transaction.Tx) {
	serializedData, err := tx.Serialize()

	if err != nil {
		fmt.Printf("Transaction serialization error: %s\n", err)
		return
	}
	data := GobEncode(Tx{
		AddrFrom:    nodeAddress,
		Transaction: serializedData,
	})

	data = append(CommandToBytes("gettx"), data...)
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
func sendInv(addr string, kind MESSAGE_TYPE, inventories [][]byte) {
	data := GobEncode(Inv{
		AddrFrom: nodeAddress,
		Type:     kind,
		data:     inventories,
	})

	data = append(CommandToBytes("inv"), data...)
	sendData(addr, data)
}

/*
handle functions receives all the info and you guessed it handles all the encoded streams of data
*/

// receives all the address from the network
func HandleAddress(request []byte) {
	//send address sends all the known nodes address, now we have to decode it
	var buff bytes.Buffer
	var payload Address

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)

	if err != nil {
		log.Panic(err)
	}

	knownNodes = append(knownNodes, payload.AddrList...)

	for _, node := range knownNodes {
		// request blocks with all the nodes that we have recieved
		SendGetBlocks(node)
	}
}

// adds the received block to the chain
func HandleBlock(request []byte, bChain *blockchain.BlockChain) {
	var buff bytes.Buffer
	var payload Block

	buff.Write(request[commandLength:])
	err := gob.NewDecoder(&buff).Decode(&payload)

	if err != nil {
		log.Panic(err)
	}

	// TODO: Implement Add block to blockchain method
	// block, err := blockchain.UnmarshalJSONTOBlock(payload.Block)
	// bChain.AddToBlockchain(block)
}

// response to get block request
func HandleGetBlock(request []byte, chain *blockchain.BlockChain) {
	var buff bytes.Buffer
	var payload GetBlocks

	buff.Write(request[commandLength:])
	err := gob.NewDecoder(&buff).Decode(&payload)

	if err != nil {
		log.Panic(err)
	}

	blocks := chain.GetBlockHashes()

	sendInv(payload.AddrFrom, BLOCK_TYPE, blocks)
}

func HandleGetData(request []byte, chain *blockchain.BlockChain) {
	var buff bytes.Buffer
	var payload GetData

	buff.Write(request[commandLength:])
	err := gob.NewDecoder(&buff).Decode(&payload)

	if err != nil {
		log.Panic(err)
	}

	if payload.Type == BLOCK_TYPE {
		block, err := chain.GetBlock([]byte(payload.data))

		if err != nil {
			return
		}

		SendBlock(payload.AddrFrom, &block)
	}

	if payload.Type == TX_TYPE {
		txId := hex.EncodeToString(payload.data)
		tx := memoryPool[txId]

		sendTx(payload.AddrFrom, tx)
	}
}

func HandleVersion(request []byte, chain *blockchain.BlockChain) {
	var buff bytes.Buffer
	var payload Version

	buff.Write(request[commandLength:])
	err := gob.NewDecoder(&buff).Decode(&payload)

	if err != nil {
		log.Panic(err)
	}

	bestHeight := chain.GetHeight()
	otherheight := payload.Height

	// if the best height is less than the height on the network then request get blocks
	if bestHeight < otherheight {
		SendGetBlocks(payload.AddressFrom)
	} else if bestHeight > otherheight {
		SendVersion(payload.AddressFrom, chain)
	}

	if !contains(knownNodes, payload.AddressFrom) {
		knownNodes = append(knownNodes, payload.AddressFrom)
	}
}

func HandleTx(request []byte, chain *blockchain.BlockChain) {
	var buff bytes.Buffer
	var payload Tx

	buff.Write(request[commandLength:])
	err := gob.NewDecoder(&buff).Decode(&payload)

	if err != nil {
		log.Panic(err)
	}

	tx, err := transaction.Deserialize(payload.Transaction)

	if err != nil {
		return
	}

	txHash := tx.CalculateHash()
	memoryPool[hex.EncodeToString(tx.CalculateHash())] = *tx

	if nodeAddress == knownNodes[0] {
		for _, node := range knownNodes {
			if node != nodeAddress && node != payload.AddrFrom {
				sendInv(node, TX_TYPE, [][]byte{txHash})
			}
		}
	} else {
		// TODO: Fix the number of nodes to mine
		if len(memoryPool) >= 2 {
			//TODO: Mine Transaction
			// MineTx(tx)
		}

	}
}

func HandleInv(request []byte) {
	var buff bytes.Buffer
	var payload Inv

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)

	// what this does is reads from the buff buffer and stores the decoded info in payload struct
	// also could have done err := gob.NewDecoder(&buff).Decode(&payload)
	err := dec.Decode(&payload)

	if err != nil {
		log.Panic(err)
	}

	// TODO: Maybe create a string from the type data for just information
	log.Printf("Received %d inventories of type %d", len(payload.data), payload.Type)

	if payload.Type == BLOCK_TYPE {
		blocksInTransit := payload.data

		blockHash := payload.data[0]
		sendGetData(payload.AddrFrom, BLOCK_TYPE, blockHash)

		// TODO: Check this out for more details(what this code does)
		newInTransit := [][]byte{}
		for _, b := range blocksInTransit {
			if bytes.Compare(b, blockHash) != 0 {
				newInTransit = append(newInTransit, b)
			}
		}
		blocksInTransit = newInTransit
	}

	if payload.Type == TX_TYPE {
		txID := payload.data[0]
		tx, _ := memoryPool[hex.EncodeToString(txID)].Serialize()
		if tx == nil {
			sendGetData(payload.AddrFrom, TX_TYPE, txID)
		}
	}

}

func HandleConnection(conn net.Conn, chain *blockchain.BlockChain) {
	// Reader is interface with read method
	req, err := ioutil.ReadAll(conn)

	defer conn.Close()
	if err != nil {
		log.Panic(err)
	}

	// get the required command
	// each request's first 12 characters is a command and rest is the load
	command := BytesToCommand(req[:12])

	switch command {
	default:
		fmt.Println("Unknown command")
		return

	case "inv":
		fmt.Println("Sending inventory")
		HandleInv(req)

	case "getversion":
		fmt.Println("Sending version")
		HandleVersion(req, chain)

	case "getdata":
		fmt.Println("Sending data of a type")
		HandleGetData(req, chain)
		break

	case "gettx":
		fmt.Println("Sending Transaction")
		HandleTx(req, chain)
		break

	case "address":
		fmt.Println("Sending known addresses")
		HandleAddress(req)
		break

	case "block":
		fmt.Println("Sending a block")
		HandleBlock(req, chain)
		break

	case "getblocks":
		HandleGetBlock(req, chain)
		break

	}
}

func contains(array []string, val string) bool {
	for _, elem := range array {
		if elem == val {
			return true
		}
	}

	return false
}

// Gob Encode
// Details: https://pkg.go.dev/encoding/gob
func GobEncode(data interface{}) []byte {
	var buff bytes.Buffer

	// the encoded data is stored in buff and the data to be encoded is `data`
	err := gob.NewEncoder(&buff).Encode(data)

	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
