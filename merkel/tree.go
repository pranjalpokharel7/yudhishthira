package merkel

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"strconv"

	"github.com/pranjalpokharel7/yudhishthira/blockchain"
)

// node struct, to encompass data
type Node struct {
	HashValue    []byte `json:"hash"` // contains the hashed byte
	parent       *Node  // parent node
	right        *Node
	left         *Node
	Tx           blockchain.Tx `json:"tx"` // transaction for data storage
	hashStrategy func(blockchain.Tx) hash.Hash
	tree         *MerkelTree
}

// Merkel tree to store all the info
type MerkelTree struct {
	Root         *Node   `json:"rootNode"`
	LeafNodes    []*Node `json:"leafNodes"`
	hashStrategy func([]byte) string
}

// hash transaction struct
func hashTransaction(tx blockchain.Tx) []byte {
	data := []byte(strconv.Itoa(tx.InputCount))
	hash := sha256.Sum256(data)
	return hash[:]
}

func hashDataSha256(data []byte) string {
	hash := sha256.Sum256(data)
	return string(hash[:])
}

func CreateMerkelTree(transactions []blockchain.Tx, tree *MerkelTree) (*MerkelTree, error) {
	if len(transactions) == 0 {
		return nil, errors.New("Can't create a new tree from empty list")
	}

	if tree == nil {
		tree = &MerkelTree{
			hashStrategy: hashDataSha256,
		}
	}

	// add to the roots of the merkel tree
	for _, tx := range transactions {
		node := Node{
			HashValue: hashTransaction(tx),
			parent:    nil,
			right:     nil,
			left:      nil,
			Tx:        tx,
		}

		tree.LeafNodes = append(tree.LeafNodes, &node)
	}

	var err error
	tree.Root, err = createMerkelTreeIntermediate(tree.LeafNodes, tree)

	return tree, err
}

func AddDataMerkelTree(tree *MerkelTree, transactions ...blockchain.Tx) (*MerkelTree, error) {
	for _, tx := range transactions {
		node := &Node{
			Tx:        tx,
			parent:    nil,
			HashValue: hashTransaction(tx),
		}
		tree.LeafNodes = append(tree.LeafNodes, node)
	}

	var err error
	tree.Root, err = createMerkelTreeIntermediate(tree.LeafNodes, tree)

	return tree, err
}

// creates a merkel tree with root specified
func createMerkelTreeIntermediate(nodes []*Node, tree *MerkelTree) (*Node, error) {
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

		contentHash := append(nodes[left].HashValue, nodes[right].HashValue...)
		hash := tree.hashStrategy(contentHash)

		n := &Node{
			left:      nodes[left],
			right:     nodes[right],
			HashValue: []byte(hash),
			tree:      tree,
		}

		nodeList = append(nodeList, n)

		nodes[left].parent = n
		nodes[right].parent = n
		if len(nodes) == 2 {
			return n, nil
		}
	}

	return createMerkelTreeIntermediate(nodeList, tree)
}

func (node *Node) Print() {
	if node == nil {
		return
	}
	fmt.Printf("%x\n", node.HashValue)
	node.left.Print()
	node.right.Print()

}

func (tree *MerkelTree) GetRoot() *Node {
	return tree.Root
}

func (tree *MerkelTree) GetLengthLeaves() int {
	return len(tree.LeafNodes)
}

func (tree *MerkelTree) VerifyTransaction(tx blockchain.Tx) bool {
	for _, node := range tree.LeafNodes {
		if bytes.Compare(hashTransaction(tx), node.HashValue) == 0 {
			parentNode := node.parent
			for parentNode != nil {
				rightHash := parentNode.right.HashValue
				leftHash := parentNode.left.HashValue

				if bytes.Compare(parentNode.HashValue, []byte(tree.hashStrategy(append(leftHash, rightHash...)))) != 0 {
					return false
				}

				parentNode = parentNode.parent
			}

			return true
		}
	}

	return false
}

func (tree MerkelTree) MarshalToJSON() ([]byte, error) {
	tree_json, err := json.Marshal(tree)
	return tree_json, err
}

func UnMarshalFromJSON(jsonData []byte) (*MerkelTree, error) {
	var unmarshalInterface map[string]interface{}
	if err := json.Unmarshal(jsonData, &unmarshalInterface); err != nil {
		fmt.Println(err)
		return nil, err
	}

	var tree *MerkelTree = &MerkelTree{}

	for k, v := range unmarshalInterface {
		if k == "rootNode" {
			tree.Root = &Node{}
			tree.Root = HandleNodeValue(v.(map[string]interface{}))
		}

		if k == "leafNodes" {
			tree.LeafNodes = HandleNodesArray(v.([]interface{}))
		}
	}

	// complete the merkel tree from leaf nodes
	var err error
	tree.Root, err = createMerkelTreeIntermediate(tree.LeafNodes, tree)

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
			node.Tx = *HandleTransaction(v.(map[string]interface{}))
		}
	}

	return node
}

func HandleTransaction(jsonData map[string]interface{}) *blockchain.Tx {
	tx := &blockchain.Tx{}

	for k, v := range jsonData {
		if k == "inputCount" {
			tx.InputCount = int(v.(float64))
		} else if k == "outputCount" {
			tx.OutputCount = int(v.(float64))
		} else if k == "itemHash" {
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
		}
	}

	return tx
}
