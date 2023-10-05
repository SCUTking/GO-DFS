package slaves

import (
	"GO-DFS/config"
	"GO-DFS/grpcPb"
	"GO-DFS/model"
	"GO-DFS/utils"
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// 修剪路径中特殊的字符
func (c *SlaveServer) TrimFileNameSpecialChar(str string) string {
	// trim special char in filename for example: #@%$^&*()_+{}|:"<>?[];',./
	reg := regexp.MustCompile(`[\/\\\:\*\?\"\<\>\|\(\)\[\]\{\}\,\;\'\#\@\%\&\$\^\~\=\+\-\!\~\s]`)
	return strings.Replace(reg.ReplaceAllString(str, ""), "...", "", -1)
}

func (s *SlaveServer) SaveUploadFile(data []byte, fileInfo *model.FileInfo) (*model.FileInfo, error) {
	var (
		err     error
		outFile *os.File
		folder  string
		fi      os.FileInfo
	)
	_, fileInfo.Name = filepath.Split(fileInfo.Name)

	// trim special char in filename for example: #@%$^&*()_+{}|:"<>?[];',./
	fileInfo.Name = s.TrimFileNameSpecialChar(fileInfo.Name)
	if s.Conf.RenameFile {
		//path.Ext(fileInfo.Name)获取后缀
		fileInfo.ReName = s.FileUtil.Hash(utils.Uuid()) + path.Ext(fileInfo.Name)
	}

	//生成文件目录的路径 folder
	folder = time.Now().Format("20060102/15/04") //根据这个格式
	//if Config().PeerId != "" {
	//	folder = fmt.Sprintf(folder+"/%s", Config().PeerId)
	//}

	if fileInfo.User != "" {
		if fileInfo.Scene != "" {
			folder = fmt.Sprintf(config.STORE_DIR_NAME+"/%s/%s/%s", fileInfo.User, fileInfo.Scene, folder)
		} else {
			folder = fmt.Sprintf(config.STORE_DIR_NAME+"/%s/%s", fileInfo.User, folder)
		}
	} else {
		return fileInfo, errors.New("User is need")
	}

	//检查文件是否存在
	if !s.FileUtil.FileExists(folder) {
		if err = os.MkdirAll(folder, 0775); err != nil {
			s.Logger.Error(err.Error())
		}
	}
	outPath := fmt.Sprintf(folder+"/%s", fileInfo.Name)
	if fileInfo.ReName != "" {
		outPath = fmt.Sprintf(folder+"/%s", fileInfo.ReName)
	}
	if s.FileUtil.FileExists(outPath) {
		if s.Conf.EnableDistinctFile {
			//循环直到有空闲的目录
			for i := 0; i < 10000; i++ {
				outPath = fmt.Sprintf(folder+"/%d_%s", i, filepath.Base(fileInfo.Name))
				fileInfo.Name = fmt.Sprintf("%d_%s", i, fileInfo.Name)
				if !s.FileUtil.FileExists(outPath) {
					break
				}
			}
		} else {

		}
	}

	s.Logger.Info(fmt.Sprintf("upload: %s", outPath))
	//由操作系统创建文件
	if outFile, err = os.Create(outPath); err != nil {
		return fileInfo, err
	}
	defer outFile.Close()
	//将利用io函数file中的内容复制进去
	reader := bytes.NewReader(data)
	if _, err = io.Copy(outFile, reader); err != nil {
		s.Logger.Error(err.Error())
		return fileInfo, errors.New("(error)fail," + err.Error())
	}

	//获取文件的大小信息
	if fi, err = outFile.Stat(); err != nil {
		s.Logger.Error(err.Error())
		return fileInfo, errors.New("(error)fail," + err.Error())
	} else {
		fileInfo.Size = fi.Size()
	}

	//补充FileInfo的其他信息
	v := "" // c.util.GetFileSum(outFile, Config().FileSumArithmetic)

	if s.Conf.EnableDistinctFile {
		v = s.FileUtil.GetFileSum(outFile, s.Conf.FileSumArithmetic)
	} else {

		v = s.FileUtil.GetFileSum(outFile, s.Conf.FileSumArithmetic)
	}
	fileInfo.Md5 = v

	//fileInfo.Path = folder //strings.Replace( folder,DOCKER_DIR,"",1)
	//路径的拼凑
	//fileInfo.Path = strings.Replace(folder, DOCKER_DIR, "", 1)
	//fileInfo.Path = folder
	//TODO 添加存储该文件host到peer
	//fileInfo.Peers = append(fileInfo.Peers, c.host)
	//fmt.Println("upload", fileInfo)
	return fileInfo, nil

}

func (s *SlaveServer) SaveUploadFileWithEC(data []byte, fileInfo *model.FileInfo) (*model.FileInfo, error) {
	var (
		err     error
		outFile *os.File
		folder  string
		fi      os.FileInfo
	)
	_, fileInfo.Name = filepath.Split(fileInfo.Name)

	// trim special char in filename for example: #@%$^&*()_+{}|:"<>?[];',./
	fileInfo.Name = s.TrimFileNameSpecialChar(fileInfo.Name)
	if s.Conf.RenameFile {
		//path.Ext(fileInfo.Name)获取后缀
		fileInfo.ReName = s.FileUtil.Hash(utils.Uuid()) + path.Ext(fileInfo.Name)
	}

	//生成文件目录的路径 folder
	folder = time.Now().Format("20060102/15/04") //根据这个格式
	//if Config().PeerId != "" {
	//	folder = fmt.Sprintf(folder+"/%s", Config().PeerId)
	//}

	if fileInfo.User != "" {
		if fileInfo.Scene != "" {
			folder = fmt.Sprintf(config.STORE_DIR_NAME+"/%s/%s/%s", fileInfo.User, fileInfo.Scene, folder)
		} else {
			folder = fmt.Sprintf(config.STORE_DIR_NAME+"/%s/%s", fileInfo.User, folder)
		}
	} else {
		return fileInfo, errors.New("User is need")
	}

	//检查文件目录是否存在
	if !s.FileUtil.FileExists(folder) {
		if err = os.MkdirAll(folder, 0775); err != nil {
			s.Logger.Error(err.Error())
		}
	}

	//分片存储
	s.Encoder(data, fileInfo.ReName, folder, fileInfo)

	//创建一个临时文件 用于生成Hash值 并计算文件大小
	outFile, err = os.CreateTemp("", "temp*")
	defer os.Remove(outFile.Name()) // clean up
	reader := bytes.NewReader(data)
	if _, err = io.Copy(outFile, reader); err != nil {
		s.Logger.Error(err.Error())
		return fileInfo, errors.New("(error)fail," + err.Error())
	}

	//获取文件的大小信息
	if fi, err = outFile.Stat(); err != nil {
		s.Logger.Error(err.Error())
		return fileInfo, errors.New("(error)fail," + err.Error())
	} else {
		fileInfo.Size = fi.Size()
	}

	//补充FileInfo的其他信息
	v := "" // c.util.GetFileSum(outFile, Config().FileSumArithmetic)

	if s.Conf.EnableDistinctFile {
		v = s.FileUtil.GetFileSum(outFile, s.Conf.FileSumArithmetic)
	} else {

		v = s.FileUtil.GetFileSum(outFile, s.Conf.FileSumArithmetic)
	}
	fileInfo.Md5 = v

	//fileInfo.Path = folder
	//TODO 添加存储该文件host到peer

	return fileInfo, nil

}

func (s *SlaveServer) saveShardsToNode(nodeAddr model.Addr, req grpcPb.ShardsReq) string {
	conn, err := grpc.Dial(utils.AddrToString(nodeAddr), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		s.Logger.Error("did not connect:" + err.Error())
	}
	defer conn.Close()
	client := grpcPb.NewGreeterClient(conn)

	// 执行RPC调用并打印收到的响应数据
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	shards, err := client.SaveShards(ctx, &req)

	s.Logger.Error(shards.Replay)
	return shards.Replay

}
