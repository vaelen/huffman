package huffman

import (
	//"fmt"
	"io"
	"math"
	"strconv"
)

const BlockSize = 65535

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
