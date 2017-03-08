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

package main

import (
	"fmt"
	"os"
	"flag"
	"path/filepath"
	"strings"
	"github.com/vaelen/huffman"
)

func OpenFiles(infile string, outfile string) (*os.File, *os.File, error) {
	input, err := os.Open(infile)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	output, err := os.Create(outfile)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}	

	return input, output, nil
}

func Encode(infile string, outfile string) error {
	fmt.Printf("Encoding: %s -> %s\n", infile, outfile)
	input, output, err := OpenFiles(infile, outfile)
	defer input.Close()
	defer output.Close()
	if err != nil {
		return err
	}
	err = huffman.Encode(input, output)
	if err != nil {
		return err
	}
	inStats, _ := input.Stat()
	outStats, _ := output.Stat()
	if inStats != nil && outStats != nil {
		inSize := float64(inStats.Size())
		outSize := float64(outStats.Size())
		compressionRatio := (outSize / inSize) * 100.0
		fmt.Printf("Compression Ratio: %2.2f%%\n", compressionRatio)
	}
	return nil
}

func Decode(infile string, outfile string) error {
	fmt.Printf("Decoding: %s -> %s\n", infile, outfile)
	input, output, err := OpenFiles(infile, outfile)
	defer input.Close()
	defer output.Close()
	if err != nil {
		return err
	}
	return huffman.Decode(input, output)
}

func main() {

	var infile string
	var outfile string
	var performDecode bool
	var showHelp bool

	flag.StringVar(&infile, "i", "", "Input file name")
	flag.StringVar(&outfile, "o", "", "Output file name")
	flag.BoolVar(&performDecode, "d", false, "Decode rather than encode")
	flag.BoolVar(&showHelp, "h", false, "Show this help message")

	flag.Parse()

	if infile == "" {
		if flag.NArg() > 0 {
			infile = flag.Arg(0)
		} else {
			showHelp = true
		}
	}

	if outfile == "" && infile != ""  {
		if flag.NArg() > 1 {
			outfile = flag.Arg(2)
		} else if performDecode {
			// Remove the "huf" extension if it exists
			p := filepath.Dir(infile)
			b := filepath.Base(infile)
			e := filepath.Ext(b)
			fmt.Println(p,b,e)
			if strings.ToLower(e) == ".huf" {
				fn := strings.TrimSuffix(b, e)
				outfile = filepath.Join(p, fn)
			} else {
				outfile = fmt.Sprintf("%s.out", infile)
			}
		} else {
			outfile = fmt.Sprintf("%s.huf", infile)
		}
	}

	if outfile != "" {
		i := 0
		p := filepath.Dir(outfile)
		b := filepath.Base(outfile)
		fn := outfile
		for _, err := os.Stat(fn); err == nil ; _, err = os.Stat(fn) {
			// File exists
			fmt.Printf("File exists: %s\n", fn)
			fn = filepath.Join(p, fmt.Sprintf("%s.%d", b, i))
			i++
		}
		outfile = fn
	}

	if showHelp {
		fmt.Fprintln(os.Stderr, "Huffman encoder/decoder. Copyright 2017, Andrew Young <andrew@vaelen.org>")
		flag.Usage()
	} else if performDecode {
		err := Decode(infile, outfile)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		err := Encode(infile, outfile)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	
}
