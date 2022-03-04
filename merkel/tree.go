package merkel

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/pranjalpokharel7/yudhishthira/blockchain"
)

// node struct, to encompass data
type Node struct {
	HashValue []byte `json:"hash"` // contains the hashed byte
	parent    *Node  // parent node
	right     *Node
	left      *Node
	Tx        blockchain.Tx `json:"tx"` // transaction for data storage
	// hashStrategy func(blockchain.Tx) hash.Hash
	tree *MerkelTree
}

// Merkel tree to store all the info
type MerkelTree struct {
	Root         *Node   `json:"rootNode"`
	LeafNodes    []*Node `json:"leafNodes"`
	hashStrategy func([]byte) []byte
}

func hashDataSha256(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

func CreateMerkelTree(transactions []blockchain.Tx, tree *MerkelTree) (*MerkelTree, error) {
	if len(transactions) == 0 {
		return nil, errors.New("can't create a new tree from empty list")
	}

	if tree == nil {
		tree = &MerkelTree{
			hashStrategy: hashDataSha256,
		}
	}

	// add to the roots of the merkel tree
	for _, tx := range transactions {
		node := Node{
			HashValue: tx.TxID,
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
			HashValue: tx.TxID,
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
	size := len(tree.LeafNodes)

	for i := 0; i < size; i++ {
		node := tree.LeafNodes[i]
		if bytes.Equal(tx.TxID, node.HashValue) { // TODO: might re-calculate transaction hash here?
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
	HashValue string        `json:"hash"` // contains the hashed byte
	Tx        blockchain.Tx `json:"tx"`   // transaction for data storage
}

func (tree MerkelTree) MarshalToJSON() ([]byte, error) {
	var treeJson TreeJson
	treeJson.Root = &NodeJson{
		HashValue: hex.EncodeToString(tree.Root.HashValue),
		Tx:        tree.Root.Tx,
	}

	for _, node := range tree.LeafNodes {
		treeJson.LeafNodes = append(treeJson.LeafNodes, &NodeJson{
			HashValue: hex.EncodeToString(node.HashValue),
			Tx:        node.Tx,
		})
	}

	tree_json, err := json.MarshalIndent(treeJson, "", "\t")
	// fmt.Println(string(tree_json))
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
