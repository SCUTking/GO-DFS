package config

import "go.uber.org/zap"

type GlobalConfig struct {
	//日志对象
	Logger *zap.Logger
	//日志处理
	Release  bool   `flagUsage:"release level of logs"`
	Silent   bool   `flagUsage:"do not print logs"`
	Log      string `json:"log" `
	LogLevel string `yaml:"logLevel" json:"logLevel"`
	//管理者
	Master string `json:"master",default:"localhost:6060"`
	//本机调试
	Debug bool `json:"debug" yaml:"debug"`
	//实际保存数据的
	Slaves           []string `json:"slaves"`
	EnableHttps      bool     `json:"enable_https"`
	Group            string   `json:"group"`
	RenameFile       bool     `json:"rename_file" yaml:"renameFile"`
	ShowDir          bool     `json:"show_dir"`
	Extensions       []string `json:"extensions"`
	RefreshInterval  int      `json:"refresh_interval"`
	EnableWebUpload  bool     `json:"enable_web_upload"`
	DownloadDomain   string   `json:"download_domain"`
	EnableCustomPath bool     `json:"enable_custom_path"`
	Scenes           []string `json:"scenes"`
	AlarmReceivers   []string `json:"alarm_receivers"`
	DefaultScene     string   `json:"default_scene"`
	//Mail                 Mail     `json:"mail"`
	AlarmUrl             string   `json:"alarm_url"`
	DownloadUseToken     bool     `json:"download_use_token"`
	DownloadTokenExpire  int      `json:"download_token_expire"`
	QueueSize            int      `json:"queue_size"`
	AutoRepair           bool     `json:"auto_repair"`
	Host                 string   `json:"host"`
	FileSumArithmetic    string   `json:"file_sum_arithmetic" yaml:"fileSumArithmetic"`
	PeerId               string   `json:"peer_id"`
	SupportGroupManage   bool     `json:"support_group_manage"`
	AdminIps             []string `json:"admin_ips"`
	AdminKey             string   `json:"admin_key"`
	EnableMergeSmallFile bool     `json:"enable_merge_small_file"`
	EnableMigrate        bool     `json:"enable_migrate"`
	EnableDistinctFile   bool     `json:"enable_distinct_file" yaml:"enableDistinctFile"`
	ReadOnly             bool     `json:"read_only"`
	EnableCrossOrigin    bool     `json:"enable_cross_origin"`
	EnableGoogleAuth     bool     `json:"enable_google_auth"`
	AuthUrl              string   `json:"auth_url"`
	EnableDownloadAuth   bool     `json:"enable_download_auth"`
	DefaultDownload      bool     `json:"default_download"`
	EnableTus            bool     `json:"enable_tus"`
	SyncTimeout          int64    `json:"sync_timeout"`
	EnableFsNotify       bool     `json:"enable_fsnotify"`
	EnableDiskCache      bool     `json:"enable_disk_cache"`
	ConnectTimeout       bool     `json:"connect_timeout"`
	ReadTimeout          int      `json:"read_timeout"`
	WriteTimeout         int      `json:"write_timeout"`
	IdleTimeout          int      `json:"idle_timeout"`
	ReadHeaderTimeout    int      `json:"read_header_timeout"`
	SyncWorker           int      `json:"sync_worker"`
	UploadWorker         int      `json:"upload_worker"`
	UploadQueueSize      int      `json:"upload_queue_size"`
	RetryCount           int      `json:"retry_count"`
	SyncDelay            int64    `json:"sync_delay"`
	WatchChanSize        int      `json:"watch_chan_size"`
	ImageMaxWidth        int      `json:"image_max_width"`
	ImageMaxHeight       int      `json:"image_max_height"`
	//Proxies              []Proxy  `json:"proxies"`
	EnablePprofDebug bool `json:"enable_pprof_debug"`

	//redis  元数据存储
	RedisAddr     string `json:"redisAddr" yaml:"redisAddr"`
	RedisPassword string `json:"redisPassword" yaml:"redisPassword"`
	//最大空闲连接数
	RedisMaxIdle int `json:"redis_max_idle" default:"30"`
	//最大连接数
	RedisMaxActive int `json:"redisMaxActive" yaml:"redisMaxActive" default:"30"`
	//空闲连接保活时间
	RedisIdleTimeout int `yaml:"redisIdleTimeout" json:"redisIdleTimeout" default:"200"`

	//zookeeper
	ZookeeperAddr string `yaml:"zookeeperAddr" json:"zookeeperAddr"`
	ManagerPath   string `yaml:"managerPath" json:"managerPath"`
	SlavesPath    string `yaml:"slavesPath" json:"slavesPath"`
	Interval      int    `yaml:"interval" json:"interval"`
	//MongoDB
	MongoDBAddr string `json:"mongoDBAddr" yaml:"mongoDBAddr"`
}
