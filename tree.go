// Huffman encoder/decoder
// Copyright 2017, Andrew Young <andrew@vaelen.org>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
//    This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
//     You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

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
