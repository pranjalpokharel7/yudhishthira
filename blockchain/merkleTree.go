package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/pranjalpokharel7/yudhishthira/utility"
)

// node struct, to encompass data
type Node struct {
	HashValue   utility.HexByte `json:"hash"` // contains the hashed byte
	parent      *Node           // parent node
	right       *Node
	left        *Node
	Transaction Tx `json:"tx"` // transaction for data storage
	tree        *MerkleTree
}

// Merkle tree to store all the info
type MerkleTree struct {
	Root         *Node   `json:"rootNode"`
	LeafNodes    []*Node `json:"leafNodes"`
	hashStrategy func([]byte) []byte
}

func hashDataSha256(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

func CreateMerkleTree(transactions []Tx, tree *MerkleTree) (*MerkleTree, error) {
	if len(transactions) == 0 {
		return nil, errors.New("can't create a new tree from empty list")
	}

	if tree == nil {
		tree = &MerkleTree{
			hashStrategy: hashDataSha256,
		}
	}

	// add to the roots of the merkle tree
	for _, tx := range transactions {
		// get hash of the transaction
		hashTx, _ := tx.CalculateTxHash()
		node := Node{
			HashValue:   hashTx,
			parent:      nil,
			right:       nil,
			left:        nil,
			Transaction: tx,
		}

		tree.LeafNodes = append(tree.LeafNodes, &node)
	}

	var err error
	tree.Root, err = createMerkleTreeIntermediate(tree.LeafNodes, tree)

	return tree, err
}

func AddDataMerkleTree(tree *MerkleTree, transactions ...Tx) (*MerkleTree, error) {
	for _, tx := range transactions {
		// get hash of the transaction
		hashTx, _ := tx.CalculateTxHash()
		node := &Node{
			Transaction: tx,
			parent:      nil,
			HashValue:   hashTx,
		}
		tree.LeafNodes = append(tree.LeafNodes, node)
	}

	var err error
	tree.Root, err = createMerkleTreeIntermediate(tree.LeafNodes, tree)

	return tree, err
}

// creates a merkle tree with root specified
func createMerkleTreeIntermediate(nodes []*Node, tree *MerkleTree) (*Node, error) {
	var nodeList []*Node

	if len(nodes) == 1 {
		return nodes[0], nil
	}

	if tree.hashStrategy == nil {
		tree.hashStrategy = hashDataSha256
	}

	for i := 0; i < len(nodes); i += 2 {
		var left, right int = i, i + 1

		if i == len(nodes)-1 {
			right = i
		}

		// get the hash of the intermediate node from left and right node
		contentHash := append(nodes[left].HashValue, nodes[right].HashValue...)
		hash := tree.hashStrategy(contentHash)

		n := &Node{
			left:      nodes[left],
			right:     nodes[right],
			HashValue: hash,
			tree:      tree,
		}

		nodeList = append(nodeList, n)

		nodes[left].parent = n
		nodes[right].parent = n
		if len(nodes) == 2 {
			return n, nil
		}
	}

	return createMerkleTreeIntermediate(nodeList, tree)
}

func (node *Node) Print() {
	if node == nil {
		return
	}
	fmt.Printf("%x\n", node.HashValue)
	node.left.Print()
	node.right.Print()
}

func (tree *MerkleTree) GetRoot() *Node {
	return tree.Root
}

func (tree *MerkleTree) GetLengthLeaves() int {
	return len(tree.LeafNodes)
}

func (tree *MerkleTree) VerifyTransaction(tx Tx) bool {
	size := len(tree.LeafNodes)

	for i := 0; i < size; i++ {
		node := tree.LeafNodes[i]
		var hashTx []byte

		//TODO: Check this
		if tx.TxID == nil {
			hashTx = tx.TxID
		} else {
			var err error
			hashTx, err = tx.CalculateTxHash()

			utility.ErrThenLogFatal(err)
		}
		if bytes.Equal(hashTx, node.HashValue) {
			parentNode := node.parent
			for parentNode != nil {
				rightHash := parentNode.right.HashValue
				leftHash := parentNode.left.HashValue

				if !bytes.Equal(parentNode.HashValue, []byte(tree.hashStrategy(append(leftHash, rightHash...)))) {
					return false
				}

				parentNode = parentNode.parent
			}

			return true
		}
	}

	return false
}

type TreeJson struct {
	Root      *NodeJson   `json:"rootNode"`
	LeafNodes []*NodeJson `json:"leafNodes"`
}

type NodeJson struct {
	HashValue   string `json:"hash"` // contains the hashed byte
	Transaction Tx     `json:"tx"`   // transaction for data storage
}

func (tree MerkleTree) MarshalToJSON() ([]byte, error) {
	var treeJson TreeJson
	treeJson.Root = &NodeJson{
		HashValue:   hex.EncodeToString(tree.Root.HashValue),
		Transaction: tree.Root.Transaction,
	}

	for _, node := range tree.LeafNodes {
		treeJson.LeafNodes = append(treeJson.LeafNodes, &NodeJson{
			HashValue:   hex.EncodeToString(node.HashValue),
			Transaction: node.Transaction,
		})
	}

	tree_json, err := json.MarshalIndent(treeJson, "", "\t")
	return tree_json, err
}

// TODO: Unmarshaling doesn't work for root node. Empty tx value for root node. Code doesn't make sense.
func UnmarshalMerkleFromInterface(unmarshalInterface map[string]interface{}) (*MerkleTree, error) {

	// var unmarshalInterface map[string]interface{}
	// if err := json.Unmarshal(jsonData, &unmarshalInterface); err != nil {
	// 	fmt.Println(err)
	// 	return nil, err
	// }

	var tree *MerkleTree = &MerkleTree{}

	for k, v := range unmarshalInterface {
		if k == "rootNode" {
			tree.Root = &Node{}
			tree.Root = HandleNodeValue(v.(map[string]interface{}))
		}

		if k == "leafNodes" {
			tree.LeafNodes = HandleNodesArray(v.([]interface{}))
		}
	}

	// complete the merkle tree from leaf nodes
	var err error
	tree.Root, err = createMerkleTreeIntermediate(tree.LeafNodes, tree)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return tree, nil
}

func HandleNodesArray(nodeArray []interface{}) []*Node {
	var nodes []*Node

	for _, nodeJsonData := range nodeArray {
		node := HandleNodeValue(nodeJsonData.(map[string]interface{}))
		nodes = append(nodes, node)
	}

	return nodes
}

func HandleNodeValue(jsonData map[string]interface{}) *Node {

	node := &Node{}

	for k, v := range jsonData {
		if k == "hash" {
			node.HashValue = []byte(v.(string))
		} else if k == "tx" {
			node.Transaction = *HandleTransaction(v.(map[string]interface{}))
		}
	}

	return node
}

func HandleTransaction(jsonData map[string]interface{}) *Tx {
	tx := &Tx{}

	for k, v := range jsonData {
		if k == "itemHash" {
			if v != nil {
				tx.ItemHash = v.([]byte)
			}
		} else if k == "sellerHash" {
			if v != nil {
				tx.SellerHash = v.([]byte)
			}
		} else if k == "buyerHash" {
			if v != nil {
				tx.BuyerHash = v.([]byte)
			}
		} else if k == "amount" {
			tx.Amount = uint64(v.(float64))
		} else if k == "txID" {
			if v != nil {
				tx.TxID = v.([]byte)
			}
		} else if k == "UTXOID" {
			if v != nil {
				tx.UTXOID = v.([]byte)
			}
		} else if k == "signature" {
			if v != nil {
				tx.Signature = v.([]byte)
			}
		}
	}

	return tx
}
