package slaves

//
//import (
//	"GO-DFS/grpcPb"
//	"fmt"
//	"golang.org/x/net/context"
//	"google.golang.org/grpc"
//	"net"
//)
//
//// hello server
//
//type server struct {
//	grpcPb.UnimplementedGreeterServer
//}
//
//func (s *server) SayHello(ctx context.Context, in *grpcPb.HelloRequest) (*grpcPb.HelloResponse, error) {
//	return &grpcPb.HelloResponse{Replay: "Hello " + in.Name}, nil
//}
//
//func main() {
//	// 监听本地的 8972 端口
//	lis, err := net.Listen("tcp", ":8972")
//	if err != nil {
//		fmt.Printf("failed to listen: %v", err)
//		return
//	}
//	s := grpc.NewServer()                      // 创建gRPC服务器
//	grpcPb.RegisterGreeterServer(s, &server{}) // 在gRPC服务端注册服务
//	// 启动服务
//	err = s.Serve(lis)
//	if err != nil {
//		fmt.Printf("failed to serve: %v", err)
//		return
//	}
//}
