package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/robfig/config"
)

var currentFilePath string // 程序的运行目录
func main() {
	fmt.Println(time.Now().UTC().Format("2006-01-02_15:04:05"), "当前版本: 2.5, 服务开启成功...")
	// 获取程序的配置
	currentFilePath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	files, _ := os.ReadDir(currentFilePath + "/config")

	fmt.Println("运行地址" + currentFilePath + "/config")
	ch := make(chan int, len(files))
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".conf") {
			continue
		}

		configSubProcess, _ := config.ReadDefault(currentFilePath + "/config/" + file.Name())
		go startBackInterval(configSubProcess, ch)
	}

	for i := 0; i < len(files); i++ {
		<-ch
	}
	fmt.Println("进程结束")
}

func startBackInterval(config *config.Config, ch chan int) {
	user, _ := config.String("", "user")
	pwd, _ := config.String("", "pwd")
	host, _ := config.String("", "host")
	port, _ := config.String("", "port")
	savedir, _ := config.String("", "savedir")
	savedir = currentFilePath + savedir
	dbsStr, _ := config.String("", "dbs")
	maxfiles, _ := config.Int("", "maxfiles")
	if maxfiles <= 0 {
		maxfiles = 180
	}
	dbs := strings.Split(dbsStr, ",")

	// 备份时间
	backtime, err := config.Float("", "backtime")

	fmt.Println("数据库正在备份准备中:", host, "。保存路径为:"+savedir, "最大保存数量：", maxfiles)
	// 如果是不float类型 则是时间判断
	if err != nil {
		backTimeStr, _ := config.String("", "backtime")
		hourMinuts := strings.Split(backTimeStr, ":")
		housr, _ := strconv.Atoi(hourMinuts[0])
		minuts, _ := strconv.Atoi(hourMinuts[1])
		for true {
			now := time.Now().UTC()
			if now.Hour() == housr && now.Minute() == minuts {
				invokeBack(user, pwd, host, port, savedir, dbs, maxfiles)
			}
			time.Sleep(time.Second * time.Duration(40))
		}
	} else {
		for true {
			time.Sleep(time.Minute * time.Duration(backtime))
			invokeBack(user, pwd, host, port, savedir, dbs, maxfiles)
		}
	}
	ch <- 1
}

func invokeBack(user string, pwd string, host string, port string, savedir string, dbs []string, maxfiles int) {
	// 判断文件数量是否大于要求值
	files, _ := os.ReadDir(savedir)
	for len(files) > maxfiles {
		minCreateTimeFile := getMinModifyTimeFile(savedir)
		if minCreateTimeFile != nil {
			os.Remove(path.Join(savedir, minCreateTimeFile.Name()))
		}
		files, _ = os.ReadDir(savedir) // 重新获取文件列表
	}

	//  mysqldump --single-transaction --column-statistics=0 --host www.52hhx.com --port 3307 -uroot -proot8114359 --databases vmq > /usr/local/backmysql/blog/vmq_2024-05-12_14:55:41.sql
	for _, db := range dbs {
		// 文件名字
		// sqlFileName := db + "_" + time.Now().UTC().Format("2006-01-02_15:04:05") + ".sql"
		// 将上面的赋值改为下面的赋值
		sqlFileNamePath := fmt.Sprintf("%s%s_%s.sql", savedir, db, time.Now().UTC().Format("2006-01-02_15:04:05"))
		backShell := fmt.Sprintf("mysqldump --single-transaction --host %s --port %s -u%s -p%s --databases %s > %s",
			host, port, user, pwd, db, sqlFileNamePath)
		fmt.Println("备份命令", backShell)
		// 创建目录 区分不同平台
		var shellOrCmd string
		if runtime.GOOS == "windows" {
			shellOrCmd = "cmd.exe"
		} else {
			shellOrCmd = "bash"
		}
		retMkdir := exec.Command(shellOrCmd, "-c", "mkdir -p "+savedir)

		retMkdirBytes, err := retMkdir.Output()
		if err != nil {
			fmt.Println("创建目录 出现错误", string(retMkdirBytes), err.Error())
		}

		retFrp := exec.Command(shellOrCmd, "-c", backShell)
		retFrpBytes, err := retFrp.CombinedOutput() // 获取标准输出和错误输出

		if err != nil {
			fmt.Println("出现错误", string(retFrpBytes), err.Error())
			fmt.Println("重新执行一次备份")
			// 先删除旧的文件
			os.Remove(sqlFileNamePath)
			invokeBack(user, pwd, host, port, savedir, []string{db}, maxfiles)
		} else {
			fmt.Println("数据库 ", db, " 备份完毕")
		}
	}
}

// 获取创建时间最长的文件
func getMinModifyTimeFile(path string) os.FileInfo {
	files, _ := os.ReadDir(path)
	if len(files) <= 0 {
		return nil
	}
	var minModifyTimeFile os.FileInfo // 创建时间最小的文件信息

	for i := 0; i < len(files); i++ {
		curFile, err := files[i].Info()
		if err != nil {
			fmt.Println("获取文件信息出错", err.Error())
		}
		if minModifyTimeFile == nil || minModifyTimeFile.ModTime().UnixNano() > curFile.ModTime().UnixNano() {
			minModifyTimeFile = curFile
		}
	}

	return minModifyTimeFile
}
