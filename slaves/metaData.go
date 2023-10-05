package slaves

import (
	"GO-DFS/model"
	"GO-DFS/utils/redis"
)

func (s *SlaveServer) SaveFileMetaData(fileInfo *model.FileInfo) {

	err := redis.SetForever(fileInfo.Md5, fileInfo)

	if err != nil {
		s.Logger.Error(err.Error())
	}

}
