package slaves

import (
	"GO-DFS/grpcPb"
	"GO-DFS/model"
	"fmt"
	"github.com/klauspost/reedsolomon"
	"os"
	"strconv"
)

const (
	dataShards = 2 // 数据分片数
	parShards  = 1 // 校验分片数
)

func (s *SlaveServer) Encoder(data []byte, fileName string, dir string, fileInfo *model.FileInfo) *model.FileInfo {

	// 创建编码矩阵
	enc, err := reedsolomon.New(dataShards, parShards)
	checkErr(err)

	//checkErr(err)

	// Split the file into equally sized shards.
	// 将数据分成相同的等分 用于Encode函数
	shards, err := enc.Split(data)
	checkErr(err)
	fmt.Printf("File split into %d data+parity shards with %d bytes/shard.\n", len(shards), len(shards[0]))

	// Encode parity
	// 分片 encode 下 将分片生成校验
	err = enc.Encode(shards)

	checkErr(err)

	// Write out the resulting files.

	//TODO saveShardsToNode函数的使用  呢

	var peers []string
	var shardNames []string
	slavesWithHttp := s.GetSlavesWithHttp()
	for i, shard := range shards {
		//加上.i的后缀
		outfn := fmt.Sprintf("%s.%d", fileName, i)
		shardGrpc := grpcPb.ShardsReq{Folder: dir, FileName: outfn, Content: shard}
		slave := chooseSlave(i, slavesWithHttp)
		//获取真正的文件名
		saveFileName := s.saveShardsToNode(slave, shardGrpc)
		var url = dir + "/" + saveFileName
		shardNames = append(shardNames, url)
		peers = append(peers, slave.Ip)

	}

	fileInfo.Peers = peers
	fileInfo.ShardNames = shardNames

	return fileInfo

	//for i := range shards {
	//	outfn := fmt.Sprintf("%s.%d", fileName, i)
	//	fmt.Println("Creating", outfn)
	//	//保存文件的路径要存在
	//	_, err := os.Create(filepath.Join(dir, outfn))
	//	checkErr(err)
	//}
	//
	//for i, shard := range shards {
	//	outfn := fmt.Sprintf("%s.%d", fileName, i)
	//
	//	fmt.Println("Writing to", outfn)
	//	// 分割的分片写入各个文件中 0777 权限
	//	err = os.WriteFile(filepath.Join(dir, outfn), shard, os.ModePerm)
	//	checkErr(err)
	//}
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		os.Exit(2)
	}
}

func Decode() {
	// Create matrix
	enc, err := reedsolomon.New(dataShards, parShards)
	checkErr(err)

	fname := "./tmp/encoder/rs.txt"
	// Create shards and load the data.
	shards := make([][]byte, dataShards+parShards)
	//iNum用于记录丢失的分片，用于分片的数据恢复
	var iNum []int
	for i := range shards {
		infn := fmt.Sprintf("%s.%d", fname, i)
		fmt.Println("Opening", infn)
		shards[i], err = os.ReadFile(infn)
		if err != nil {
			fmt.Println("Error reading file", err)
			shards[i] = nil
			iNum = append(iNum, i)
		}
	}

	// Verify the shards  校验各个分片是否正常 能否构建出一个文件
	ok, err := enc.Verify(shards)

	if ok {
		fmt.Println("No reconstruction needed")
	} else {
		fmt.Println("Verification failed. Reconstructing data")
		// 如果分片还足以恢复切片，就恢复分片，但是完整性不能得到验证，还得调用enc.Verify检查
		err = enc.Reconstruct(shards)
		if err != nil {
			fmt.Println("Reconstruct failed -", err)
			os.Exit(1)
		}
		// 再次验证分片
		ok, err = enc.Verify(shards)
		if !ok {
			fmt.Println("Verification failed after reconstruction, data likely corrupted.")
			os.Exit(1)
		}
		checkErr(err)
		// 丢失的分片数据恢复
		for _, i2 := range iNum {
			outfn := fname + "." + strconv.Itoa(i2)
			err = os.WriteFile(outfn, shards[i2], os.ModePerm)

		}
	}

	outfn := "./tmp/decoder/rs.txt"
	fmt.Println("Writing data to", outfn)
	f, err := os.Create(outfn)
	checkErr(err)

	// We don't know the exact filesize.
	// 还原数据
	err = enc.Join(f, shards, len(shards[0])*dataShards)
	checkErr(err)
}

// TODO 从服务器选择的优化
func chooseSlave(i int, slavesWithHttp []model.Addr) model.Addr {
	// i只能是三个副本  所以i取值为0到2

	if len(slavesWithHttp) == 1 {
		return slavesWithHttp[0]
	} else if len(slavesWithHttp) == 2 {
		if i == 0 {
			return slavesWithHttp[0]
		} else {
			return slavesWithHttp[1]
		}
	} else if len(slavesWithHttp) == 0 {
		return model.Addr{}
	} else {
		//大于三个的时候就返回前几个即可
		return slavesWithHttp[i]
	}
}
