package merkel

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"hash"
	"strconv"

	"github.com/pranjalpokharel7/yudhishthira/transaction"
)

type Node struct {
	hashValue    string // contains the hashed byte
	parent       *Node  // parent node
	right        *Node
	left         *Node
	tx           transaction.Tx // transaction for data storage
	hashStrategy func(transaction.Tx) hash.Hash
	tree         *MerkelTree
}

func hashTransaction(tx transaction.Tx) []byte {
	data := []byte(strconv.Itoa(tx.InputCount))
	hash := sha1.Sum(data)
	return hash[:]
}

type MerkelTree struct {
	root *Node
	// this function is the hashing function so that we can test with multiple hashing functions, alsow not fixed how to hash a transaction, so a place holder
	hashStrategy func([]byte) []byte
}

func CreateMerkelTree(transactions []transaction.Tx, tree *MerkelTree) (*MerkelTree, error) {
	if len(transactions) == 0 {
		return nil, errors.New("Can't create a new tree from empty list")
	}

	if tree == nil {
		tree = &MerkelTree{}
	}

	var nodes []*Node

	// add to the roots of the merkel tree
	for _, tx := range transactions {
		node := Node{
			hashValue: string(hashTransaction(tx)),
			parent:    nil,
			right:     nil,
			left:      nil,
			tx:        tx,
		}

		nodes = append(nodes, &node)
	}

	var err error
	tree.root, err = createMerkelTreeIntermediate(nodes, tree)

	return tree, err
}

// creates a merkel tree with root specified
func createMerkelTreeIntermediate(nodes []*Node, tree *MerkelTree) (*Node, error) {
	var nodeList []*Node

	if len(nodes) == 1 {
		return nodes[0], nil
	}

	for i := 0; i < len(nodes); i += 2 {
		var left, right int = i, i + 1

		if i == len(nodes)-1 {
			right = i
		}

		contentHash := []byte(nodes[left].hashValue + nodes[right].hashValue)
		hash := sha1.Sum(contentHash)

		n := &Node{
			left:      nodes[left],
			right:     nodes[right],
			hashValue: string(hash[:]),
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
	fmt.Printf("%x ", node.hashValue)
	node.left.Print()
	node.right.Print()

}

func (tree *MerkelTree) GetRoot() *Node {
	return tree.root
}
