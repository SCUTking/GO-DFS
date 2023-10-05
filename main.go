package main

import (
	"GO-DFS/config"
	"GO-DFS/grpcPb"
	"GO-DFS/manager"
	"GO-DFS/model"
	"GO-DFS/slaves"
	"GO-DFS/utils"
	"GO-DFS/utils/redis"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/koding/multiconfig"
	"github.com/samuel/go-zookeeper/zk"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	g errgroup.Group
)

//type server struct {
//	grpcPb.UnimplementedGreeterServer
//}
//
//func (s *server) SayHello(ctx context.Context, in *grpcPb.HelloRequest) (*grpcPb.HelloResponse, error) {
//	return &grpcPb.HelloResponse{Replay: "Hello " + in.Name}, nil
//}

func main() {

	//初始化
	conf := loadConf()
	redis.Setup(conf)
	////
	////initManagerServer(conf)
	//initAllServerByConf(conf)

	ServiceRegistry(conf)

	if err := g.Wait(); err != nil {
		conf.Logger.Error(err.Error())
	}

	//fmt.Println(config)
	//conn, err2 := redis.Dial("tcp", "43.140.198.154:6379")
	//conn.Do("AUTH", "123456")
	//_, err41 := conn.Do("set", "1", "2")
	//if err41 != nil {
	//	fmt.Println(err41)
	//}

	//err2 := redisUtil.Set("1", "1", 1)
	//if err2 != nil {
	//	fmt.Println(err2)
	//}

}

func initLogger(conf *config.GlobalConfig) (logger *zap.Logger) {
	if conf.Silent {
		logger = zap.NewNop()
		return
	}

	var err error
	if conf.Release {
		logger, err = zap.NewProduction()
	} else {

		encoderConfig := zap.NewProductionEncoderConfig()
		//时间格式
		encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
		}
		jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)

		file, _ := os.OpenFile(conf.Log, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		sync := zapcore.AddSync(file)
		var core zapcore.Core
		if conf.LogLevel == "INFO" {
			core = zapcore.NewCore(jsonEncoder, sync, zapcore.InfoLevel)
		} else if conf.LogLevel == "DEBUG" {
			core = zapcore.NewCore(jsonEncoder, sync, zapcore.DebugLevel)
		}

		logger = zap.New(core)

	}
	if err != nil {
		log.Fatalln("init logger failed ", err)
	}
	return
}

func loadConf() *config.GlobalConfig {
	m := multiconfig.NewWithPath("config/config.yaml")
	config := new(config.GlobalConfig)
	err := m.Load(config)
	//初始化日志类
	config.Logger = initLogger(config)

	if err != nil {
		config.Logger.Error("Failed to load configuration")
	}
	fmt.Println(config)
	return config
}

func managerRegistry(conn *zk.Conn, config *config.GlobalConfig, server *manager.Server) {

	managerStr := utils.AddrToString(server.Master)
	managerData := []byte(managerStr)
	_, err := conn.Create(config.ManagerPath, managerData, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err != nil {
		config.Logger.Error("Central server registration for zookeeper failed" + "because of :" + err.Error())
	}
	fmt.Println("Registered manager server")

	go func() {
		for {
			//清空之前的记录
			server.Slaves = server.Slaves[:0]

			// 获取从服务器节点数据
			children, _, err := conn.Children(config.SlavesPath)
			for i := 0; i < len(children); i++ {
				var path string
				path = config.SlavesPath + "/" + children[i]
				data, _, err := conn.Get(path)
				if err != nil {
					config.Logger.Error("Failed to get secondary server data")
				} else {
					config.Logger.Info("Secondary server data:" + string(data))
					addr := utils.StringToAddr(string(data))
					server.Slaves = append(server.Slaves, addr)
				}
			}
			if err != nil {
				config.Logger.Error("Failed to get secondary server data")
			}
			// 可以在这里添加更多的逻辑来判断从服务器状态

			time.Sleep(time.Second * time.Duration(config.Interval)) // 5秒钟后再次检测
		}
	}()

}

func slavesRegistry(conn *zk.Conn, config *config.GlobalConfig, ip string, port string) {

	slavesData := []byte(ip + ":" + port)
	path := config.SlavesPath + "/" + utils.Uuid()
	_, err := conn.Create(path, slavesData, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err != nil {
		config.Logger.Error("Failed to register zookeeper from slavesServer")
	}
	fmt.Println("Registered slaves server")

}

func ServiceRegistry(config *config.GlobalConfig) {

	zookeeperAddr := config.ZookeeperAddr
	fmt.Println("地址" + zookeeperAddr)
	zkServers := []string{zookeeperAddr}
	conn, _, err := zk.Connect(zkServers, time.Second*5)
	if err != nil {
		config.Logger.Error("Failed to connect to zookeeper:" + zookeeperAddr)
	}

	//创建一个从服务器的持久化节点  防止后面注册从服务器的时候找不到开始的目录导致报错
	isExist, _, err := conn.Exists(config.SlavesPath)
	if !isExist {
		_, err = conn.Create(config.SlavesPath, nil, 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	server := manager.NewServer(config)

	//两个只能使用一个  所以不需要再对zookeeper的连接进行处理
	initSlavesServers(config, conn)
	initManagerServer(config, conn, server)

}

func initSlavesServers(config *config.GlobalConfig, conn *zk.Conn) {
	slavesArr := config.Slaves
	var addrs []model.Addr
	//将slaves的String数组转化为Addr
	for i := 0; i < len(slavesArr); i++ {
		split := strings.Split(slavesArr[i], ":")
		atoi, err := strconv.Atoi(split[1])
		if err != nil {
			config.Logger.Error("Type conversion failed")
		}
		var a model.Addr
		a.Ip = split[0]
		a.Port = atoi
		addrs = append(addrs, a)
	}

	//说明要在本机启动对应的端口服务
	var ports []int
	for _, v := range addrs {
		local, err := utils.CheckIfLocal(v.Ip)
		if err != nil {
			config.Logger.Error("Check Ip failed")
		}
		if local {
			ports = append(ports, v.Port)
		}
	}

	ip, _ := utils.ExternalIP()
	for i := 0; i < len(ports); i++ {
		//server := &http.Server{
		//	Addr:         ":" + strconv.Itoa(ports[i]),
		//	Handler:      initSlavesMux(s),
		//	ReadTimeout:  5 * time.Second,
		//	WriteTimeout: 10 * time.Second,
		//}

		// 监听本地的 端口
		lis, err := net.Listen("tcp", ":"+strconv.Itoa(ports[i]))
		s := grpc.NewServer() // 创建gRPC服务器
		server := &slaves.SlaveServer{
			FileUtil:   new(utils.FileUtil),
			Logger:     config.Logger,
			ResultUtil: new(utils.ResultUtil),
			Conf:       config,
		}

		grpcPb.RegisterGreeterServer(s, server) // 在gRPC服务端注册服务

		//注册从服务器

		slavesRegistry(conn, config, ip.String(), strconv.Itoa(ports[i]))

		// 启动服务
		if config.Debug {
			//不阻塞进程 允许本机开启主和从节点
			go func() {
				err = s.Serve(lis)
				if err != nil {
					fmt.Printf("failed to serve: %v", err)
					return
				}
			}()
		} else {
			err = s.Serve(lis)
			if err != nil {
				fmt.Printf("failed to serve: %v", err)
				return
			}
		}

		//g.Go(func() error {
		//	return server.ListenAndServe()
		//})
	}

}
func initManagerServer(config *config.GlobalConfig, conn *zk.Conn, s *manager.Server) {
	master := config.Master

	masterAddr := utils.StringToAddr(master)
	local, err := utils.CheckIfLocal(masterAddr.Ip)

	//将主机的IP转为公网IP
	ip, _ := utils.ExternalIP()
	addr := utils.StringToAddr(ip.String() + ":" + strconv.Itoa(masterAddr.Port))
	s.Master = addr

	if err != nil {
		config.Logger.Error("Check Ip failed")
	}
	if local {
		server := &http.Server{
			Addr:         ":" + strconv.Itoa(masterAddr.Port),
			Handler:      initManagerMux(s),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		}
		managerRegistry(conn, config, s)

		g.Go(func() error {
			return server.ListenAndServe()
		})
	}

}

func initManagerMux(server *manager.Server) http.Handler {
	r := gin.New()
	//设置文件上传的限制
	r.MaxMultipartMemory = 1 << 20
	r.GET("/getAliveSlavesInfo", server.GetAliveSlavesInfo)
	r.POST("/upload", server.Upload)
	r.GET("/getAllUser", server.GetAllUser)
	r.POST("/addUser", server.AddUser)
	r.POST("/addUserScene", server.AddUserScene)
	r.POST("/test", server.SayHello)
	return r
}

func initSlavesMux(server *manager.Server) http.Handler {
	r := gin.New()
	r.GET("/", func(c *gin.Context) {
		var result model.JsonResult
		result.Data = "test1"
		result.Status = "success1"
		result.Message = "test1"
		c.JSON(http.StatusOK, result)
	})
	return r
}
