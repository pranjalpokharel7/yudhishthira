package merkel

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
	"strconv"

	"github.com/pranjalpokharel7/yudhishthira/transaction"
)

// node struct, to encompass data
type Node struct {
	HashValue    []byte // contains the hashed byte
	parent       *Node  // parent node
	right        *Node
	left         *Node
	tx           transaction.Tx // transaction for data storage
	hashStrategy func(transaction.Tx) hash.Hash
	tree         *MerkelTree
}

// Merkel tree to store all the info
type MerkelTree struct {
	root         *Node
	leafNodes    []*Node
	hashStrategy func([]byte) string
}

// hash transaction struct
func hashTransaction(tx transaction.Tx) []byte {
	data := []byte(strconv.Itoa(tx.InputCount))
	hash := sha256.Sum256(data)
	return hash[:]
}

func hashDataSha256(data []byte) string {
	hash := sha256.Sum256(data)
	return string(hash[:])
}

func CreateMerkelTree(transactions []transaction.Tx, tree *MerkelTree) (*MerkelTree, error) {
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
			tx:        tx,
		}

		tree.leafNodes = append(tree.leafNodes, &node)
	}

	var err error
	tree.root, err = createMerkelTreeIntermediate(tree.leafNodes, tree)

	return tree, err
}

func AddDataMerkelTree(tree *MerkelTree, transactions ...transaction.Tx) (*MerkelTree, error) {
	for _, tx := range transactions {
		node := &Node{
			tx:        tx,
			parent:    nil,
			HashValue: hashTransaction(tx),
		}
		tree.leafNodes = append(tree.leafNodes, node)
	}

	var err error
	tree.root, err = createMerkelTreeIntermediate(tree.leafNodes, tree)

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
	return tree.root
}

func (tree *MerkelTree) GetLengthLeaves() int {
	return len(tree.leafNodes)
}

func (tree *MerkelTree) VerifyTransaction(tx transaction.Tx) bool {
	for _, node := range tree.leafNodes {
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
