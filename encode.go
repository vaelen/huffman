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
	"bufio"
	"io"
	"math"
	"strconv"
	"container/heap"
	"encoding/binary"
)

const BlockSize = 65535

func Encode(input io.Reader, output io.Writer) error {
	b := make([]byte, BlockSize)
	defer fmt.Println()
	for {
		n, err := input.Read(b)
		if err != nil && err != io.EOF {
			return err
		}
		if n > 0 {
			// fmt.Printf("Processing %d bytes\n", n)
			fmt.Print("#")
			e2 := EncodeChunk(b[:n], output)
			if e2 != nil {
				return e2
			}
		}
		if err == io.EOF {
			break
		}
	}
	return nil
}

func EncodeChunk(input []byte, output io.Writer) error {
	o := bufio.NewWriter(output)
	defer o.Flush()
	tree := buildEncodingTree(input)
	
	m := buildEncodingMap(make(map[byte]string), "", tree)

	writeChunkHeader(tree, uint16(len(input)), output)
	
	buf := make([]byte,0)
	var err error
	
	for _,v := range input {
		bits := []byte(m[v])
		for _, b := range bits {
			buf = append(buf, b)
		}
		buf, err = encodeBytes(buf, o)
		if err != nil {
			return err
		}
	}
	if len(buf) > 0 {
		buf = padToByte(buf, '0')
		encodeBytes(buf, o)
	}
	if err != nil {
		return err
	}
	return nil
}

func writeChunkHeader(tree *TreeNode, dataSize uint16, output io.Writer) error {
	header := buildHeader(make([]byte, 0), tree)
	headerSizeBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(headerSizeBytes, uint16(len(header)))
	dataSizeBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(dataSizeBytes, dataSize)
	// fmt.Printf("Header Size: %d bytes\n", len(header))
	// fmt.Printf("Data Size: %d bytes\n", dataSize)
	output.Write(headerSizeBytes)
	output.Write(header)
	output.Write(dataSizeBytes)
	return nil
}

func buildHeader(header []byte, i *TreeNode) []byte {
	if i == nil {
		return header
	} else if i.IsLeaf() {
		return append(header, 1, i.value)
	} else {
		header = append(header, 0)
		header = buildHeader(header, i.left)
		header = buildHeader(header, i.right)
		return header
		
	}
}	

func buildEncodingTree(input []byte) *TreeNode {
	counts := make(map[byte]int)
	for _, v := range input {
		c, ok := counts[v]
		if (!ok) {
			c = 0
		}
		counts[v] = c + 1
	}
	
	q := make(PriorityQueue, 0)
	i := 0
	for k, v := range counts {
		heap.Push(&q, &TreeNode{value: k, priority: v})
		i++
	}

	heap.Init(&q)

	for q.Len() > 1 {
		leftTreeNode := heap.Pop(&q).(*TreeNode)
		rightTreeNode := heap.Pop(&q).(*TreeNode)

		newTreeNode := &TreeNode{
			priority: (leftTreeNode.priority + rightTreeNode.priority),
			left: leftTreeNode,
			right: rightTreeNode,
		}

		heap.Push(&q, newTreeNode)
	}

	return heap.Pop(&q).(*TreeNode)

}

func buildEncodingMap(m map[byte]string, prefix string, i *TreeNode) map[byte]string {
	if i != nil {
		if i.IsLeaf() {
			m[i.value] = prefix
		} else {
			buildEncodingMap(m, prefix + "0", i.left)
			buildEncodingMap(m, prefix + "1", i.right)
		}
	}
	return m
}

func padToByte(buf []byte, pad byte) []byte {
	if len(buf) > 0 {
		for math.Mod(float64(len(buf)), 8.0) > 0.0 {
			// Pad with zeros
			buf = append(buf, pad)
		}
	}
	return buf
}

func encodeBytes(bits []byte, output io.ByteWriter) ([]byte, error) {
	for len(bits) > 7 {
		s := string(bits[:8])
		bits = bits[8:]
		b, err := strconv.ParseUint(s, 2, 8)
		if err != nil {
			return bits, err
		}
		output.WriteByte(byte(b))
	}
	return bits, nil
}
