package huffman

import (
	"fmt"
	"errors"
	"bufio"
	"io"
	"encoding/binary"
)

func Decode(input io.Reader, output io.Writer) error {
	i := bufio.NewReader(input)
	o := bufio.NewWriter(output)
	defer o.Flush()
	for {
		tree, size, err := readChunkHeader(i)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		m := buildDecodingMap(make(map[string]byte), "", tree)

		fmt.Printf("Processing %d original bytes\n", size)
		err = DecodeChunk(i, o, m, size)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
	}
	return nil
}

func DecodeChunk(input io.ByteReader, output io.ByteWriter, m map[string]byte, size uint16) error {
	return nil
}

func readChunkHeader(input io.Reader) (*TreeNode, uint16, error) {
	headerSizeBytes := make([]byte, 2)
	n, err := input.Read(headerSizeBytes)
	if n < 2 || err == io.EOF {
		return nil, 0, io.EOF
	} else if err != nil {
		return nil, 0, err
	}
	headerSize := binary.BigEndian.Uint16(headerSizeBytes)

	var tree *TreeNode
	var in io.ByteReader
	in, ok := input.(io.ByteReader)
	if !ok {
		in = bufio.NewReader(input)
	}
	tree, _, err = readHeader(in, headerSize)
	if err != io.EOF {
		return nil, 0, err
	}
	printTree("", tree)

	
	dataSizeBytes := make([]byte, 2)
	n, err = input.Read(dataSizeBytes)
	if n < 2 || err == io.EOF {
		return nil, 0, io.EOF
	} else if err != nil {
		return nil, 0, err
	}
	dataSize := binary.BigEndian.Uint16(headerSizeBytes)

	return tree, dataSize, nil
}

func readHeader(input io.ByteReader, bytesLeft uint16) (*TreeNode, uint16, error) {
	if input == nil || bytesLeft == 0 {
		return nil, 0, nil
	}
	node := TreeNode{}
	t, err := input.ReadByte()
	if err != nil {
		return nil, 0, err
	}
	bytesLeft--
	switch t {
	case 0:
		// Branch
		node.left, bytesLeft, err = readHeader(input, bytesLeft)
		node.right, bytesLeft, err = readHeader(input, bytesLeft)
	case 1:
		// Leaf
		node.value, err = input.ReadByte()
		if err != nil {
			return nil, 0, err
		}
	default:
		return nil, 0, errors.New(fmt.Sprintf("Invalid Node Type: %d", t))
	}
	return &node, bytesLeft, nil
}

func buildDecodingMap(m map[string]byte, prefix string, i *TreeNode) map[string]byte {
	if i != nil {
		if i.IsLeaf() {
			m[prefix] = i.value
		} else {
			buildDecodingMap(m, prefix + "0", i.left)
			buildDecodingMap(m, prefix + "1", i.right)
		}
	}
	return m
}
