package manager

import (
	"GO-DFS/model"
	"math/rand"
	"time"
)

//负载均衡策略的选择
//根据从服务器的网络情况，cpu占用等方面选择从服务器

func (s *Server) ChooseSlave() model.Addr {
	slaves := s.Slaves

	//随机选取
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(slaves))
	return slaves[i]
}
