package manager

import (
	"GO-DFS/config"
	"GO-DFS/model"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/net/context"
)

func (s *Server) CheckUser(username string, password string) bool {

	b := s.IfExistUser(username)
	if b {
		//存在 验证密码是否正确
	} else {
		return false
	}
	return true

}

func (s *Server) IfExistUser(username string) bool {
	var result model.User
	collection := s.MgDBClient.Database(config.MongoDBDataBase).Collection("user")
	filter := bson.D{{"name", username}}
	err := collection.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		s.Logger.Error("user not exist")
		return false
	}
	return true
}

// 根据判断数据库中有没有goDfs数据库来  判断是否初始化MongoDB
func (s *Server) IfInitMonGODB() bool {
	databaseNames, err := s.MgDBClient.ListDatabaseNames(context.TODO(), bson.M{})
	if err != nil {
		s.Logger.Error(err.Error())
	}
	for _, name := range databaseNames {
		if name == config.MongoDBDataBase {
			return true
		}
	}
	return false
}
