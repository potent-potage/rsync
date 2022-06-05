package main

import (
	"fmt"
	"os"
	"io"

	"github.com/minio/sha256-simd"
)

func adler32(s string) uint32 {
    var a uint32 = 1
    var b uint32 = 0

	for _, ch := range s {
		a += uint32(ch)
		b += a
	}

	return ((b % m) << 16) | (a % m)
}

const (
	m = 1 << 16
	BlockSize = 1024 * 6
	// Hasher = md5.New()
	// sha256.New()
)

func fastSignature(block []byte) (uint32, uint32, uint32) {
	var a, b uint32
	l := uint32(len(block))

	for i, v := range block {
		a += uint32(v)
		b += (l - uint32(i) + uint32(1)) * uint32(v)
	}

	r1 := a % m
	r2 := b % m
	r := (r2 << 16) | r1
	return r1, r2, r
}

type BlockData struct {
	Index uint64
	FastSig uint32
	StrongSig []byte
}

func calculateDelta(signatures []BlockData) {
	signatureTable := make(map[uint32][]BlockData)
	for _, b := range signatures {
		signatureTable[b.FastSig] = append(signatureTable[b.FastSig], b)
	}

	type span struct {
		k int
		l int
	}
}

func readBlocks(filename string) []BlockData {
	// var r io.Reader

	// n, err := r.Read(buffer)
	// if err != nil && err != io.EOF {

	blocks := make([]BlockData, 0)
	i := uint64(0)

    f, err := os.Open(filename)
    if err != nil {
        panic(err)
    }
    defer f.Close()

    buf := make([]byte, 10)

    for {
        n, err := f.Read(buf)
        if err == io.EOF {
            // there is no more data to read
            break
        }
        if err != nil {
            fmt.Println(err)
            continue
        }
		if n > 0 {
			block := buf[:n]
			sha := sha256.New()
			sha.Write(block)
			strongSig := sha.Sum(nil)
			_, _, fastSig := fastSignature(block)
			bd := BlockData{
				Index:		i,
				FastSig:	fastSig,
				StrongSig:	strongSig,
			}
			blocks = append(blocks, bd)
        }
    }

	return blocks
}

func main() {
	// s := "Wikipedia"
	// x := adler32(s)
	// hex_value := fmt.Sprintf("%x", x)
	// fmt.Printf("Hex value of %d is = %s\n", x, hex_value)

	// r1, r2, r := rollingHash([]byte(s))
	// r1, r2, r := fastSignature([]byte(s))
	// hex_value = fmt.Sprintf("%x", r1)
	// fmt.Printf("Hex value of %d is = %s\n", x, hex_value)
	// hex_value = fmt.Sprintf("%x", r2)
	// fmt.Printf("Hex value of %d is = %s\n", x, hex_value)
	// hex_value = fmt.Sprintf("%x", r)
	// fmt.Printf("Hex value of %d is = %s\n", x, hex_value)

	readBlocks("test.txt")
}
