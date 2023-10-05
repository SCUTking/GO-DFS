package config

import "GO-DFS/model"

const (
	MongoDBDataBase                = "goDfs"
	STORE_DIR_NAME                 = "files"
	LOG_DIR_NAME                   = "log"
	DATA_DIR_NAME                  = "data"
	CONF_DIR_NAME                  = "conf"
	STATIC_DIR_NAME                = "static"
	CONST_STAT_FILE_COUNT_KEY      = "fileCount"
	CONST_BIG_UPLOAD_PATH_SUFFIX   = "/big/upload/"
	CONST_STAT_FILE_TOTAL_SIZE_KEY = "totalSize"
	CONST_Md5_ERROR_FILE_NAME      = "errors.md5"
	CONST_Md5_QUEUE_FILE_NAME      = "queue.md5"
	CONST_FILE_Md5_FILE_NAME       = "files.md5"
	CONST_REMOME_Md5_FILE_NAME     = "removes.md5"
	CONST_SMALL_FILE_SIZE          = 1024 * 1024
)

var (
	SCENE_ALREADY_EXIST    = model.JsonResult{Message: "将要创建的Scene已经存在", Status: "fail", Data: nil}
	PARAM_TYPE_ERROR       = model.JsonResult{Message: "参数类型错误", Status: "fail", Data: nil}
	USER_ACCOUNT_NOT_EXIST = model.JsonResult{Message: "该用户不存在", Status: "fail", Data: nil}
)
