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
	"errors"
	"bufio"
	"io"
	"encoding/binary"
)

func Decode(input io.Reader, output io.Writer) error {
	i := bufio.NewReader(input)
	o := bufio.NewWriter(output)
	defer o.Flush()
	defer fmt.Println()
	for {
		tree, dataSize, err := readChunkHeader(i)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		// fmt.Printf("Processing %d original bytes\n", dataSize)
		fmt.Print("#")
		err = DecodeChunk(i, o, tree, dataSize)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
	}
	return nil
}

func DecodeChunk(input io.ByteReader, output io.ByteWriter, tree *TreeNode, dataSize uint16) error {
	buf := make([]byte, 0)
	for {
		b, err := input.ReadByte()
		if err != nil {
			return err
		}
		bitString := fmt.Sprintf("%08b", b)
		bits := []byte(bitString)
		for _, bit := range bits {
			buf = append(buf, bit)
		}
		tryAgain := true
		for tryAgain {
			tryAgain = false
			node := tree
			for i, bit := range buf {
				switch bit {
				case '0':
					node = node.left
				case '1':
					node = node.right
				default:
					// This should never happen
					return errors.New(fmt.Sprintf("Invalid Bit: %s", string(bit)))
				}
				if node == nil {
					// This should never happen
					return errors.New("Invalid Decoding Tree Found.")
				} else if node.IsLeaf() {
					// Found a match
					output.WriteByte(node.value)
					buf = buf[i + 1:]
					dataSize--
					//fmt.Printf("Found: %s, Buffer Size: %d, Bytes Left: %d\n", string(node.value), len(buf), dataSize)
					tryAgain = true
					break
				}
			}
			if dataSize == 0 {
				tryAgain = false
			}
		}
		if dataSize == 0 {
			break
		}
	}
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
	// fmt.Printf("Header Size: %d bytes\n", headerSize)
	
	var tree *TreeNode
	var in io.ByteReader
	in, ok := input.(io.ByteReader)
	if !ok {
		in = bufio.NewReader(input)
	}
	tree, _, err = readHeader(in, headerSize)
	if err != nil && err != io.EOF {
		return nil, 0, err
	}
	// printTree("", tree)
	
	dataSizeBytes := make([]byte, 2)
	n, err = input.Read(dataSizeBytes)
	if n < 2 || err == io.EOF {
		return nil, 0, io.EOF
	} else if err != nil {
		return nil, 0, err
	}
	dataSize := binary.BigEndian.Uint16(dataSizeBytes)
	// fmt.Printf("Data Size: %d bytes\n", dataSize)

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
	// fmt.Printf("Read Type Byte:  %#02x, Bytes Left: %d\n", t, bytesLeft);
	switch t {
	case 0:
		// Branch
		node.left, bytesLeft, err = readHeader(input, bytesLeft)
		if err != nil {
			return nil, 0, err
		}
		node.right, bytesLeft, err = readHeader(input, bytesLeft)
		if err != nil {
			return nil, 0, err
		}
	case 1:
		// Leaf
		node.value, err = input.ReadByte()
		bytesLeft--
		// fmt.Printf("Read Value Byte: %#x, Bytes Left: %d\n", t, bytesLeft);
		if err != nil {
			return nil, 0, err
		}
	default:
		return nil, 0, errors.New(fmt.Sprintf("Invalid Node Type: %d", t))
	}
	return &node, bytesLeft, nil
}

