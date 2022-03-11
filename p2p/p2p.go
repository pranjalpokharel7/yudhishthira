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
	"github.com/pranjalpokharel7/yudhishthira/utility"
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
	// set initial knownNode
	knownNodes  = []string{"192.168.1.68:3000"} // list of all the knownNodes
	nodeAddress string                          // address of this node

	// here string is the transaction id and it point to the actual transaction
	memoryPool      = make(map[string]blockchain.Tx)
	blocksInTransit [][]byte
)

// const for types
// using integer rather than strings
// well may need to serialize this too
const (
	BLOCK_TYPE   = 1
	TX_TYPE      = 2
	VERSION_TYPE = 3
	INV_TYPE     = 4
)

type MESSAGE_TYPE int

// wrapper struct to send a block
type Block struct {
	AddrFrom string
	Block    blockchain.Block
}

// a version message that contains the version of the chain in the node
type Version struct {
	Timestamp   uint64
	AddressFrom string
	Height      uint64
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
	Data     []byte // needs to be a serialized array of bytes
}

// provides the block header hashes from a particular point
type GetBlocks struct {
	AddrFrom string
	Data     []byte
	Height   uint64
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
	Data     [][]byte     // 2D array of byte, each byte array contains a transaction
}

func CommandToBytes(cmd string) []byte {
	var bytes [commandLength]byte

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
func SendGetBlocks(addr string, chain *blockchain.BlockChain) {
	var lastHash []byte
	lastHash = chain.LastHash
	var blocks = GetBlocks{
		AddrFrom: nodeAddress,
		Data:     lastHash,
		Height:   chain.GetHeight(),
	}
	// first 12 character is command, rest is the payload
	// check this link for more details
	//READ: https://developer.bitcoin.org/reference/p2p_networking.html#headers
	info := append(CommandToBytes("getblocks"), GobEncode(blocks)...)

	sendData(addr, info)
}

// sends all the blocks
func SendBlock(addr string, block *blockchain.Block) {
	// data, err := block.MarshalBlockToJSON()
	// if err != nil {
	// 	fmt.Printf("Block serialization error: %s\n", err)
	// 	return
	// }
	var blocks = Block{
		AddrFrom: nodeAddress,
		Block:    *block,
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
		Data:     id,
	})

	data = append(CommandToBytes("getdata"), data...)
	sendData(addr, data)
}

// send a particular transaction to the given address
func sendTx(addr string, tx blockchain.Tx) {
	serializedData, err := tx.SerializeTxToGOB()

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
	inv := Inv{
		AddrFrom: nodeAddress,
		Type:     kind,
		Data:     inventories,
	}
	data := GobEncode(inv)
	var payload Inv
	gob.NewDecoder(bytes.NewBuffer(data)).Decode(&payload)
	fmt.Printf("%x\n", data)
	data = append(CommandToBytes("inv"), data...)
	sendData(addr, data)
}

/*
handle functions receives all the info and you guessed it handles all the encoded streams of data
*/

// receives all the address from the network
func HandleAddress(request []byte, chain *blockchain.BlockChain) {
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
		SendGetBlocks(node, chain)
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
	fmt.Printf("Received a block of hash: %x\n", payload.Block.BlockHash)
	bChain.AddBlock(&payload.Block)

	if len(blocksInTransit) > 0 {
		blockHash := blocksInTransit[0]
		sendGetData(payload.AddrFrom, BLOCK_TYPE, blockHash)

		blocksInTransit = blocksInTransit[1:]
	}
}

// response to get block request
func HandleGetBlocks(request []byte, chain *blockchain.BlockChain) {
	var buff bytes.Buffer
	var payload GetBlocks

	buff.Write(request[commandLength:])
	err := gob.NewDecoder(&buff).Decode(&payload)

	if err != nil {
		log.Panic(err)
	}

	// TODO: Implement get blockhashes where it gets all the blocks
	// either from particular version or entire chain hashes
	blocks := chain.GetBlockHashes(payload.Data)
	fmt.Println(string(buff.Bytes()))
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
		block, err := chain.GetBlock([]byte(payload.Data))

		if err != nil {
			log.Panic(err)
			return
		}

		SendBlock(payload.AddrFrom, block)
	}

	if payload.Type == TX_TYPE {
		txId := hex.EncodeToString(payload.Data)
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

	// height on the current chain
	bestHeight := chain.GetHeight()

	// height of received chain
	otherheight := payload.Height

	// if the best height is less than the height on the network then request get blocks
	if bestHeight < otherheight {
		fmt.Println("Sending Get block request")
		SendGetBlocks(payload.AddressFrom, chain)
	} else if bestHeight > otherheight {
		fmt.Println("Sending version of the current block")
		SendVersion(payload.AddressFrom, chain)
	} else {
		fmt.Printf("Same block height: %d", chain.GetHeight())
	}

	// if nodes are not known add them to known nodes
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

	tx, err := blockchain.DeserializeTxFromGOB(payload.Transaction)

	if err != nil {
		return
	}

	txHash, err := tx.CalculateTxHash()
	memoryPool[hex.EncodeToString(txHash)] = *tx

	if nodeAddress == knownNodes[0] {
		for _, node := range knownNodes {
			if node != nodeAddress && node != payload.AddrFrom {
				sendInv(node, TX_TYPE, [][]byte{txHash})
			}
		}
	} else {
		// TODO: Fix the number of nodes to mine
		// TODO Mine tx
		if len(memoryPool) >= 2 {

		}
	}
}

func HandleInv(request []byte) {
	buff := bytes.NewBuffer(request[commandLength:])
	var payload Inv

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(buff)

	// what this does is reads from the buff buffer and stores the decoded info in payload struct
	// also could have done err := gob.NewDecoder(&buff).Decode(&payload)
	err := dec.Decode(&payload)

	if err != nil {
		log.Panic(err)
	}

	// for printing or debugging purposes
	// can remove this one later
	// if log needs to be created then may need to use this one
	typeStringMap := map[MESSAGE_TYPE]string{
		BLOCK_TYPE:   "BLOCK",
		TX_TYPE:      "TX",
		VERSION_TYPE: "VERSION",
		INV_TYPE:     "INV",
	}
	fmt.Printf("%x\n", buff.Bytes())
	log.Printf("Received %d inventories of type %s", len(payload.Data), typeStringMap[payload.Type])

	for _, inv := range payload.Data {
		fmt.Printf("%x\n", inv)
	}

	if payload.Type == BLOCK_TYPE {
		blocksInTransit = payload.Data

		if len(payload.Data) != 0 {
			blockHash := payload.Data[0]
			sendGetData(payload.AddrFrom, BLOCK_TYPE, blockHash)

			// this section is equivalent to blocksInTransit.remove(blockHash)
			// insert all the items in the chain in blocksInTransit and the operate it later so that it gets appended to the chain
			newInTransit := [][]byte{}
			for _, b := range blocksInTransit {
				// add all the nodex except blockHash to blocksInTransit
				if bytes.Compare(b, blockHash) != 0 {
					newInTransit = append(newInTransit, b)
				}
			}
			blocksInTransit = newInTransit
		}
	}

	if payload.Type == TX_TYPE {
		txID := payload.Data[0]
		tx := memoryPool[hex.EncodeToString(txID)]
		txByte, _ := tx.SerializeTxToGOB()
		if txByte == nil {
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
	fmt.Println(command)
	switch command {
	default:
		fmt.Println("Unknown command")
		return

	case "inv":
		fmt.Println("Receiving inventory")
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
		HandleAddress(req, chain)
		break

	case "block":
		fmt.Println("Receiving a block")
		HandleBlock(req, chain)
		break

	case "getblocks":
		HandleGetBlocks(req, chain)
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

func StartServer(nodeId string) {
	nodeAddress = fmt.Sprintf("%s:%s", utility.GetNodeAddress(), nodeId)
	// minerAddress = minerAddress
	ln, err := net.Listen(protocol, nodeAddress)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	// defer chain.Database.Close()
	// go CloseDB(chain)

	chain := &blockchain.BlockChain{}
	chain = blockchain.InitBlockChain()

	if nodeAddress != knownNodes[0] {
		SendVersion(knownNodes[0], chain)
	} else {
		// chain = blockchain.InitBlockChain()
		b := blockchain.CreateBlock()
		tx := []blockchain.Tx{
			{
				Amount: 69,
			},
			{
				Amount: 6969,
			},
		}
		b.AddTransactionsToBlock(tx)
		if chain.GetHeight() == 0 {
			chain.AddBlock(b)
			chain.AddBlock(b)
			chain.AddBlock(b)
			chain.AddBlock(b)
			chain.AddBlock(b)
		}
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		go HandleConnection(conn, chain)

	}
}
