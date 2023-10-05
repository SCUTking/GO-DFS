package model

// 文件类型的API返回体
type FileResult struct {
	Url     string `json:"url"`
	Md5     string `json:"md5"`
	Path    string `json:"path"`
	Domain  string `json:"domain"`
	Scene   string `json:"scene"`
	Size    int64  `json:"size"`
	ModTime int64  `json:"mtime"`
	//Just for Compatibility
	Scenes  string `json:"scenes"`
	RetMsg  string `json:"retmsg"`
	RetCode int    `json:"retcode"`
	Src     string `json:"src"`
}

type Addr struct {
	Ip   string
	Port int
}

type JsonResult struct {
	Message string      `json:"message"`
	Status  string      `json:"status"`
	Data    interface{} `json:"data"`
}

type FileInfo struct {
	Name       string   `json:"name"`
	ReName     string   `json:"rename"`
	Path       []string `json:"path"`
	Md5        string   `json:"md5"`
	Size       int64    `json:"size"`
	Peers      []string `json:"peers"`
	ShardNames []string `json:"shardNames"`
	User       string   `json:"user"`
	Scene      string   `json:"scene"`
	TimeStamp  int64    `json:"timeStamp"`
	OffSet     int64    `json:"offset"`
	retry      int
	op         string
}

type User struct {
	Name     string  `json:"name"`
	Password string  `json:"password"`
	SceneArr []Scene `json:"sceneArr"`
	//0表示普通身份  1表示管理员身份
	Type int `json:"type"`
}

type Scene struct {
	//对应的场景
	SceneName string `json:"sceneName"`
	//对应场景对应文件权限 0表示可上传  1表示可下载 2表示两者都可以
	Type int `json:"type"`
}
