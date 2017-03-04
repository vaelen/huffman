package huffman

import (
	"fmt"
	"bufio"
	"io"
	"math"
	"container/heap"
)

func Encode(input io.Reader, output io.Writer)  {
	b := make([]byte, BlockSize)
	o := bufio.NewWriter(output)
	defer o.Flush()
	for {
		n, err := input.Read(b)
		if err != nil && err != io.EOF {
			fmt.Println(err)
			return
		}
		if n > 0 {
			fmt.Printf("Processing %d bytes\n", n)
			EncodeChunk(b[:n], o)
		}
		if err == io.EOF {
			break
		}
	}
}

func EncodeChunk(input []byte, output io.ByteWriter) {
	tree := buildEncodingTree(input)
	printTree("", tree)
	
	m := buildEncodingMap(make(map[byte]string), "", tree)

	var maxSize float64
	for _,v := range m {
		maxSize = math.Max(float64(len(v)), maxSize)
	}
	writeChunkHeader(tree, uint16(maxSize), uint16(len(input)), output)
	
	buf := make([]byte,0)
	var err error
	
	for _,v := range input {
		bits := []byte(m[v])
		for _, b := range bits {
			buf = append(buf, b)
		}
		buf, err = encodeBytes(buf, output)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	if len(buf) > 0 {
		buf = padToByte(buf, '0')
		encodeBytes(buf, output)
	}
	if err != nil {
		fmt.Println(err)
		return
	}
}

func writeChunkHeader(tree *TreeNode, bitSize uint16, dataSize uint16, output io.ByteWriter) error {
	sizeBits := fmt.Sprintf("%016b", bitSize)
	header := []byte(sizeBits)
	
	header = buildHeader(header, tree, sizeString(bitSize))
	dataSizeBits := fmt.Sprintf("%016b", dataSize);
	for _,b := range dataSizeBits {
		header = append(header, byte(b))
	}
	if len(header) > 0 {
		header = padToByte(header, '1')
		// fmt.Printf("Header Size: %d bits, %.2f bytes\n", len(header), float64(len(header))/8.0)
		_, err := encodeBytes(header, output)
		if err != nil {
			return err
		}
	}
	return nil
}

func sizeString(size uint16) string {
	return fmt.Sprintf("%%0%db", size)
}

func buildHeader(header []byte, i *TreeNode, size string) []byte {
	if i != nil {
		if i.IsLeaf() {
			bits := []byte(fmt.Sprintf(size, i.value))
			for _,b := range bits {
				header = append(header, b)
			}
		} else {
			header = append(header, '0')
			header = buildHeader(header, i.left, size)
			header = buildHeader(header, i.right, size)
		}
	}
	return header
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
