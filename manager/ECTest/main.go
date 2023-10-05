package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/klauspost/reedsolomon"
)

func main() {
	var DataShards = 10
	var ParShards = 20
	var OutDir = "./out"

	var dataShards = &DataShards
	var parShards = &ParShards
	var outDir = &OutDir

	fname := "in.txt"

	// 1.Create encoding matrix.
	enc, err := reedsolomon.NewStream(*dataShards, *parShards)
	checkErr2(err)

	fmt.Println("Opening", fname)
	f, err := os.Open(fname)
	checkErr2(err)

	instat, err := f.Stat()
	checkErr2(err)

	shards := *dataShards + *parShards
	out := make([]*os.File, shards)

	// 2.创建输入文件 30个shards
	dir, file := filepath.Split(fname)
	if *outDir != "" {
		dir = *outDir
	}
	for i := range out {
		outfn := fmt.Sprintf("%s.%d", file, i)
		fmt.Println("Creating", outfn)
		out[i], err = os.Create(filepath.Join(dir, outfn))
		checkErr2(err)
	}

	// Split into files.
	data := make([]io.Writer, *dataShards)
	for i := range data {
		data[i] = out[i]
	}
	// 3.原始文件拆分
	err = enc.Split(f, data, instat.Size())
	checkErr2(err)

	// Close and re-open the files.
	input := make([]io.Reader, *dataShards)

	for i := range data {
		out[i].Close()
		f, err := os.Open(out[i].Name())
		fmt.Println("Error ", err)
		input[i] = f
		defer f.Close()
	}

	// 4.封装 parity
	parity := make([]io.Writer, *parShards)
	for i := range parity {
		parity[i] = out[*dataShards+i]
		defer out[*dataShards+i].Close()
	}

	// 5.Encode 编码rs格式
	err = enc.Encode(input, parity)
	checkErr2(err)

	fmt.Printf("File split into %d data + %d parity shards.\n", *dataShards, *parShards)

}

func checkErr2(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
