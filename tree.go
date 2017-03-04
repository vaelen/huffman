package huffman

import (
	"fmt"
)

// A tree node
type TreeNode struct {
	value byte
	priority int
	index int
	left *TreeNode
	right *TreeNode
}

func (i *TreeNode) IsLeaf() bool {
	return i.left == nil && i.right == nil
}

func printTree(prefix string, node *TreeNode) {
	if node != nil {
		if node.IsLeaf() {
			fmt.Printf("%s-%d\n", prefix, node.value)
		} else {
			fmt.Printf("%s-*\n", prefix)
			printTree(prefix + "|", node.left)
			printTree(prefix + "|", node.right)
		}
	}
}
