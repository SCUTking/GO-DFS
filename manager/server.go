package manager

import (
	"GO-DFS/config"
	"GO-DFS/grpcPb"
	"GO-DFS/model"
	"GO-DFS/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	Master     model.Addr
	Slaves     []model.Addr
	FileUtil   *utils.FileUtil
	ResultUtil *utils.ResultUtil
	MgDBClient *mongo.Client
	Logger     *zap.Logger
	conf       *config.GlobalConfig
}

func (s *Server) GetAliveSlavesInfo(c *gin.Context) {
	c.JSON(http.StatusOK, s.Slaves)
	return
}

func (s *Server) Upload(c *gin.Context) {
	//formFile, header, _ := c.Request.FormFile("file")
	file, _ := c.FormFile("file")
	open, _ := file.Open()
	defer open.Close()
	all, _ := io.ReadAll(open)
	//log.Println(file.Filename)

	//TODO 权限控制
	user := c.PostForm("user")
	scene := c.PostForm("scene")

	// 连接到server端，此处禁用安全传输
	//TODO 选取存储节点并保存到对应的节点中去 负载均衡和远程调用从节点
	conn, err := grpc.Dial(utils.AddrToString(s.ChooseSlave()), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		s.Logger.Error("did not connect: " + err.Error())
	}
	defer conn.Close()
	client := grpcPb.NewGreeterClient(conn)

	// 执行RPC调用并打印收到的响应数据
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	uploadFileInfo := &grpcPb.UploadFileInfo{
		FileName: file.Filename,
		User:     user,
		Scene:    scene,
		Content:  all,
	}
	r, err := client.UploadFile(ctx, uploadFileInfo)
	if err != nil {
		s.Logger.Error("could not greet: " + err.Error())
	}
	s.Logger.Info("Greeting: " + r.GetReplay())

	// 上传文件至指定的完整文件路径

	c.JSON(http.StatusOK, s.ResultUtil.Successs(r))
}

func (s *Server) GetAllUser(c *gin.Context) {
	collection := s.MgDBClient.Database(config.MongoDBDataBase).Collection("user")
	cur, err := collection.Find(context.TODO(), bson.D{{}})
	defer cur.Close(context.TODO())
	if err != nil {

	}
	var results []model.User
	for cur.Next(context.TODO()) {
		//定义一个文档，将单个文档解码为result
		var result model.User
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, result)
	}
	c.JSON(http.StatusOK, results)
}

func (s *Server) AddUser(c *gin.Context) {
	user := c.Query("username")
	password := c.Query("password")
	userType := c.Query("type")
	atoi, _ := strconv.Atoi(userType)

	defaultScene := model.Scene{SceneName: "default", Type: 3}
	var sn []model.Scene
	sn = append(sn, defaultScene)
	newUser := model.User{Name: user, Password: password, Type: atoi, SceneArr: sn}
	collection := s.MgDBClient.Database(config.MongoDBDataBase).Collection("user")

	one, err := collection.InsertOne(context.TODO(), newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, false)
	}
	c.JSON(http.StatusOK, one)

}

func (s *Server) AddUserScene(c *gin.Context) {
	id := c.Query("userId")
	sceneName := c.Query("sceneName")
	sceneType := c.Query("type")

	var u model.User
	user := &u
	collection := s.MgDBClient.Database(config.MongoDBDataBase).Collection("user")
	//构造mongodb专属的ID才能进行查询
	objectID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objectID}
	findOne := collection.FindOne(context.TODO(), filter)
	if findOne.Err() != nil {
		s.Logger.Error(findOne.Err().Error())
		c.JSON(http.StatusBadRequest, config.USER_ACCOUNT_NOT_EXIST)
		return
	}

	findOne.Decode(&user)

	//重复性判断
	for _, each := range user.SceneArr {
		if sceneName == each.SceneName {
			c.JSON(http.StatusBadRequest, config.SCENE_ALREADY_EXIST)
			return
		}
	}
	var scene model.Scene
	scene.SceneName = sceneName
	var err error
	scene.Type, err = strconv.Atoi(sceneType)
	if err != nil {
		c.JSON(http.StatusBadRequest, config.PARAM_TYPE_ERROR)
		return
	}
	user.SceneArr = append(user.SceneArr, scene)

	update := bson.D{{"$set", bson.D{{"scenearr", user.SceneArr}}}}

	byID, err := collection.UpdateByID(context.TODO(), objectID, update)
	if err != nil {
		c.JSON(http.StatusOK, false)
		return
	}

	c.JSON(http.StatusOK, s.ResultUtil.Successs(byID))

}

/*
*
用于测试从服务器的接口的方法
*/
func (s *Server) SayHello(c *gin.Context) {
	//formFile, header, _ := c.Request.FormFile("file")

	conn, err := grpc.Dial(utils.AddrToString(s.ChooseSlave()), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		s.Logger.Error("did not connect: " + err.Error())
	}
	defer conn.Close()
	client := grpcPb.NewGreeterClient(conn)

	// 执行RPC调用并打印收到的响应数据
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	test := &grpcPb.HelloRequest{Name: "nihao"}
	r, err := client.SayHello(ctx, test)
	if err != nil {
		s.Logger.Error("could not greet: " + err.Error())
	}
	s.Logger.Info("Greeting: " + r.GetReplay())

	// 上传文件至指定的完整文件路径

	c.JSON(http.StatusOK, s.ResultUtil.Successs(r))
}

// 创建一个新的主服务器
func NewServer(conf *config.GlobalConfig) *Server {

	// Connect to MongoDB(初始化MonGoDB)
	clientOptions := options.Client().ApplyURI(conf.MongoDBAddr)
	var ctx = context.TODO()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		conf.Logger.Error("MongoDB connect failed")
	}

	//初始化服务器类
	server := &Server{
		FileUtil:   new(utils.FileUtil),
		MgDBClient: client,
		Logger:     conf.Logger,
		ResultUtil: new(utils.ResultUtil),
		conf:       conf,
	}

	if !server.IfInitMonGODB() {

		defaultScene := model.Scene{SceneName: "default", Type: 3}
		var s []model.Scene
		s = append(s, defaultScene)
		user := model.User{Name: "admin", Password: "admin", Type: 1, SceneArr: s}

		//if mongodb.ExistDataBase(MongoDBDataBase, conf) {
		//
		//}

		//如果没有数据库，或集合都会隐式创建
		collection := client.Database(config.MongoDBDataBase).Collection("user")

		one, err := collection.InsertOne(ctx, user)
		fmt.Println(one.InsertedID)
		if err != nil {
			server.Logger.Error(err.Error())
		}
	}

	return server
}
