package main

import (
	"container/heap"
	"fmt"
	"bufio"
	"io"
	"os"
	"math"
	"strconv"
)

const BlockSize = 65535

//// Priority Queue /////

// Based on code from https://golang.org/pkg/container/heap/


// An Item is something we manage in a priority queue.
type Item struct {
	value byte
	priority int
	index int
	left *Item
	right *Item
}

func (i *Item) IsLeaf() bool {
	return i.left == nil && i.right == nil
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

//// End Priority Queue ////


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

func padToByte(buf []byte, pad byte) []byte {
	if len(buf) > 0 {
		for math.Mod(float64(len(buf)), 8.0) > 0.0 {
			// Pad with zeros
			// fmt.Printf("Buffer Size: %d bits, %.2f bytes\n", len(buf), float64(len(buf))/8.0)
			buf = append(buf, pad)
		}
	}
	// fmt.Printf("Buffer Size: %d bits, %.2f bytes\n", len(buf), float64(len(buf))/8.0)
	return buf
}

func writeChunkHeader(tree *Item, bitSize uint16, dataSize uint16, output io.ByteWriter) error {
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

func buildHeader(header []byte, i *Item, size string) []byte {
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


func encodeBytes(bits []byte, output io.ByteWriter) ([]byte, error) {
	for len(bits) > 7 {
		s := string(bits[:8])
		bits = bits[7:]
		b, err := strconv.ParseUint(s, 2, 8)
		if err != nil {
			return bits[:8], err
		}
		output.WriteByte(byte(b))
	}
	return bits, nil
}

func buildEncodingTree(input []byte) *Item {
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
		heap.Push(&q, &Item{value: k, priority: v})
		i++
	}

	heap.Init(&q)

	for q.Len() > 1 {
		leftItem := heap.Pop(&q).(*Item)
		rightItem := heap.Pop(&q).(*Item)

		newItem := &Item{
			priority: (leftItem.priority + rightItem.priority),
			left: leftItem,
			right: rightItem,
		}

		heap.Push(&q, newItem)
	}

	return heap.Pop(&q).(*Item)

}

func buildEncodingMap(m map[byte]string, prefix string, i *Item) map[byte]string {
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

func main() {

	infile := "lorem.txt"
	outfile := "lorem.bin"
	
	input, err := os.Open(infile)
	if err != nil {
		fmt.Println(err)
		return
	}

	output, err := os.Create(outfile)
	if err != nil {
		fmt.Println(err)
		return
	}

	Encode(input, output)

}
