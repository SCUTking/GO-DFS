package slaves

import (
	"GO-DFS/model"
	"fmt"
	"github.com/goccy/go-json"
	"io"
	"net/http"
)

// GetSlavesWithHttp 从节点如何知道其他从节点服务器的地址
func (s *SlaveServer) GetSlavesWithHttp() []model.Addr {

	// 创建一个HTTP客户端
	client := &http.Client{}

	// 创建一个GET请求
	req, err := http.NewRequest("GET", "http://"+s.Conf.Master+"/getAliveSlavesInfo", nil)
	if err != nil {
		s.Logger.Error(err.Error())
	}

	// 发送请求并获取响应
	resp, err := client.Do(req)
	if err != nil {
		s.Logger.Error(err.Error())
	}
	defer resp.Body.Close()

	// 读取响应的内容
	body, err := io.ReadAll(resp.Body)
	addrs := []model.Addr{}
	//json反序列化一定要是取地址才有效
	json.Unmarshal(body, &addrs)
	fmt.Println(addrs)

	if err != nil {
		s.Logger.Error(err.Error())
	}

	// 打印响应内容
	return addrs
}
