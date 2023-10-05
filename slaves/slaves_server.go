package slaves

import (
	"GO-DFS/config"
	"GO-DFS/grpcPb"
	"GO-DFS/model"
	"GO-DFS/utils"
	"bytes"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"io"
	"os"
)

type SlaveServer struct {
	grpcPb.UnimplementedGreeterServer
	FileUtil   *utils.FileUtil
	ResultUtil *utils.ResultUtil
	Logger     *zap.Logger
	Conf       *config.GlobalConfig
}

func (s *SlaveServer) SayHello(ctx context.Context, in *grpcPb.HelloRequest) (*grpcPb.HelloResponse, error) {
	slavesWithHttp := s.GetSlavesWithHttp()
	for _, each := range slavesWithHttp {
		fmt.Println(each)
	}
	return &grpcPb.HelloResponse{Replay: "Hello " + in.Name}, nil
}

//func (s *SlaveServer) SendFile(ctx context.Context, file *grpcPb.FileTest) (*grpcPb.HelloResponse, error) {
//	info := model.FileInfo{
//		Name:  "asdasd.docx",
//		Scene: "scene",
//		User:  "user",
//	}
//
//	uploadFile, err := s.SaveUploadFile(file.GetContent(), &info)
//	return &grpcPb.HelloResponse{Replay: "Hello " + uploadFile.Md5}, err
//}

func (s *SlaveServer) UploadFile(ctx context.Context, file *grpcPb.UploadFileInfo) (*grpcPb.HelloResponse, error) {
	info := model.FileInfo{
		Name:  file.GetFileName(),
		Scene: file.GetScene(),
		User:  file.GetUser(),
	}

	uploadFile, err := s.SaveUploadFileWithEC(file.GetContent(), &info)

	fmt.Println(uploadFile)
	s.SaveFileMetaData(&info)

	return &grpcPb.HelloResponse{Replay: "Hello " + uploadFile.Md5}, err
}

func (s *SlaveServer) SaveShards(ctx context.Context, req *grpcPb.ShardsReq) (*grpcPb.HelloResponse, error) {

	var err error
	var outFile *os.File

	folder := req.GetFolder()
	//将要返回的文件名
	fileName := req.GetFileName()
	//检查文件是否存在
	if !s.FileUtil.FileExists(folder) {
		if err = os.MkdirAll(folder, 0775); err != nil {
			s.Logger.Error(err.Error())
		}
	}
	outPath := fmt.Sprintf(folder+"/%s", fileName)

	if s.FileUtil.FileExists(outPath) {
		if s.Conf.EnableDistinctFile {
			//循环直到有空闲的目录
			for i := 0; i < 10000; i++ {
				outPath = fmt.Sprintf(folder+"/%d_%s", i, fileName)
				fileName = fmt.Sprintf("%d_%s", i, fileName)
				if !s.FileUtil.FileExists(outPath) {
					break
				}
			}
		} else {

		}
	}

	//由操作系统创建文件
	if outFile, err = os.Create(outPath); err != nil {
		return &grpcPb.HelloResponse{Replay: fileName}, err
	}
	defer outFile.Close()
	//将利用io函数file中的内容复制进去
	reader := bytes.NewReader(req.GetContent())
	if _, err = io.Copy(outFile, reader); err != nil {
		s.Logger.Error(err.Error())
		return &grpcPb.HelloResponse{Replay: fileName}, errors.New("(error)fail," + err.Error())
	}

	return &grpcPb.HelloResponse{Replay: fileName}, err
}
