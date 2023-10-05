package main

import (
	"fmt"
	"io"
	"os"

	"github.com/klauspost/reedsolomon"
)

var OutFile = "out2.txt"
var outFile = &OutFile

var DataShards = 10
var ParShards = 20
var OutDir = "./out"

var dataShards = &DataShards
var parShards = &ParShards
var outDir = &OutDir

func main() {
	fname := "out/in.txt"

	// 1.Create matrix
	enc, err := reedsolomon.NewStream(*dataShards, *parShards)
	checkErr(err)

	// 2.Open the inputs
	shards, size, err := openInput(*dataShards, *parShards, fname)
	checkErr(err)

	// 3.Verify the shards
	ok, err := enc.Verify(shards)
	if ok {
		fmt.Println("No reconstruction needed")
	} else {
		fmt.Println("Verification failed. Reconstructing data")
		shards, size, err = openInput(*dataShards, *parShards, fname)
		checkErr(err)
		// 3.1 重新创建删除的文件
		out := make([]io.Writer, len(shards))
		for i := range out {
			if shards[i] == nil {
				//dir, _ := filepath.Split(fname)
				outfn := fmt.Sprintf("%s.%d", fname, i)
				fmt.Println("Creating", outfn)
				out[i], err = os.Create(outfn)
				checkErr(err)
			}
		}
		fmt.Println("reconstruct")
		// 3.2 重建30个shards
		err = enc.Reconstruct(shards, out)
		if err != nil {
			fmt.Println("Reconstruct failed -", err)
			os.Exit(1)
		}
		// Close output.
		for i := range out {
			if out[i] != nil {
				err := out[i].(*os.File).Close()
				checkErr(err)
			}
		}
		shards, size, err = openInput(*dataShards, *parShards, fname)
		ok, err = enc.Verify(shards)
		if !ok {
			fmt.Println("Verification failed after reconstruction, data likely corrupted:", err)
			os.Exit(1)
		}
		checkErr(err)
	}

	// 4.Join the shards and write them
	outfn := *outFile
	if outfn == "" {
		outfn = fname
	}

	fmt.Println("Writing data to", outfn)
	f, err := os.Create(outfn)
	checkErr(err)

	shards, size, err = openInput(*dataShards, *parShards, fname)
	checkErr(err)

	// join恢复原文件 but We don't know the exact filesize.
	err = enc.Join(f, shards, int64(*dataShards)*size)
	checkErr(err)
}

func openInput(dataShards, parShards int, fname string) (r []io.Reader, size int64, err error) {
	// Create shards and load the data.
	shards := make([]io.Reader, dataShards+parShards)
	for i := range shards {
		infn := fmt.Sprintf("%s.%d", fname, i)
		fmt.Println("Opening", infn)
		f, err := os.Open(infn)
		if err != nil {
			fmt.Println("Error reading file", err)
			shards[i] = nil
			continue
		} else {
			shards[i] = f
		}
		stat, err := f.Stat()
		checkErr(err)
		if stat.Size() > 0 {
			size = stat.Size()
		} else {
			shards[i] = nil
		}
	}
	return shards, size, nil
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
