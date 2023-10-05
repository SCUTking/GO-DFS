package utils

import (
	"GO-DFS/model"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

// 获取ip
func ExternalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected to the network?")
}

// 获取ip
func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}

func Contains(arr []string, target string) bool {
	for _, element := range arr {
		if element == target {
			return true
		}
	}
	return false
}

// 获取本机的IP地址切片
func GetLocalIps() ([]string, error) {
	var localAddr []string
	localAddr = append(localAddr, "localhost", "127.0.0.1")

	ip, err := ExternalIP()
	if err != nil {
		return localAddr, err
	}
	//默认允许在本机运行的端口
	localAddr = append(localAddr, ip.String())
	return localAddr, nil
}

// 检测IP是否是本机IP
func CheckIfLocal(ip string) (bool, error) {
	ips, err := GetLocalIps()
	if err != nil {
		return false, err
	}
	if Contains(ips, ip) {
		return true, nil
	}
	return false, nil

}

// ip加端口转化为Addr
func StringToAddr(s string) model.Addr {
	split := strings.Split(s, ":")
	var a model.Addr
	a.Ip = split[0]
	atoi, _ := strconv.Atoi(split[1])
	a.Port = atoi
	return a
}

// Addr转换为IP
func AddrToString(addr model.Addr) string {
	return addr.Ip + ":" + strconv.Itoa(addr.Port)
}

// UUID
func Uuid() string {
	newUUID := uuid.New().String()
	uuidWithoutDash := strings.ReplaceAll(newUUID, "-", "")
	return uuidWithoutDash
}

type FileUtil struct {
}

// 文件hash值的求取
func (this *FileUtil) GetFileHash(file *os.File) string {
	file.Seek(0, 0)
	hash := sha256.New()
	io.Copy(hash, file)
	sum := fmt.Sprintf("%x", hash.Sum(nil))
	return sum
}

func (this *FileUtil) Hash(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func (this *FileUtil) FileExists(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil
}

func (this *FileUtil) GetFileSum(file *os.File, alg string) string {
	alg = strings.ToLower(alg)
	if alg == "sha256" {
		return this.GetFileHash(file)
	} else {
		//TODO 添加其他的加密算法
		return this.GetFileHash(file)
	}
}

/*
*
返回工具类
*/
type ResultUtil struct {
}

func (this *ResultUtil) Successs(data interface{}) model.JsonResult {
	return this.SuccesssWithMsg(data, "")
}

func (this *ResultUtil) SuccesssWithMsg(data interface{}, msg string) model.JsonResult {
	return model.JsonResult{
		Status:  "success",
		Data:    data,
		Message: msg,
	}
}

func (this *ResultUtil) FailWithMsg(data interface{}, msg string) model.JsonResult {
	return model.JsonResult{
		Status:  "fail",
		Data:    data,
		Message: msg,
	}
}

func (this *ResultUtil) Fail(data interface{}) model.JsonResult {
	return this.FailWithMsg(data, "")
}
