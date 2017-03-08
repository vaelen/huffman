package main

import (
	"fmt"
	"os"
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
	return huffman.Encode(input, output)
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

	infile := "lorem.txt"
	outfile := "lorem.bin"
	outfile2 := "lorem.out"

	var err error

	err = Encode(infile, outfile)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = Decode(outfile, outfile2)
	if err != nil {
		fmt.Println(err)
		return
	}

}
