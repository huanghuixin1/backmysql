package main

import (
	"io/ioutil"
	"strings"
	"github.com/robfig/config"
	"time"
	"strconv"
	"os/exec"
	"fmt"
)

func main() {
	fmt.Println(time.Now().UTC().Format("2006-01-02_15:04:05"), "服务开启中...")
	// 获取程序的配置
	files, _ := ioutil.ReadDir("./config")
	ch := make(chan int, len(files))
	for _, file := range files {
		if (!strings.HasSuffix(file.Name(), ".conf")) {
			continue
		}

		configSubProcess, _ := config.ReadDefault("./config/" + file.Name())
		go startBackInterval(configSubProcess, ch)
	}

	for i := 0; i < len(files); i++ {
		<-ch
	}
}

func startBackInterval(config *config.Config,  ch chan int) {
	user, _ := config.String("", "user")
	pwd, _ := config.String("", "pwd")
	host, _ := config.String("", "host")
	port, _ := config.String("", "port")
	savedir, _ := config.String("", "savedir")
	dbsStr, _ := config.String("", "dbs")
	dbs := strings.Split(dbsStr, ",")

	backtime, err := config.Float("", "backtime")

	fmt.Println("数据库正在备份准备中:", host, "。保存路径为:" + savedir)
	// 如果是不float类型 则是时间判断
	if err != nil {
		backTimeStr, _ := config.String("", "backtime")
		hourMinuts := strings.Split(backTimeStr, ":")
		housr, _ := strconv.Atoi(hourMinuts[0])
		minuts, _ := strconv.Atoi(hourMinuts[1])
		for true {
			now := time.Now().UTC()
			if (now.Hour() == housr && now.Minute() == minuts) {
				invokeBack(user, pwd, host, port, savedir, dbs)
			}
			time.Sleep(time.Second * time.Duration(40))
		}
	} else {
		for true {
			time.Sleep(time.Minute * time.Duration(backtime))
			invokeBack(user, pwd, host, port, savedir, dbs)
		}
	}
	ch <- 1
}

func invokeBack(user string, pwd string, host string, port string, savedir string, dbs []string) {
	for _, db := range dbs {
		backShell := fmt.Sprintf("mysqldump --host %s --port %s -u%s -p%s --databases %s > %s%s.sql",
			host, port, user, pwd, db, savedir, db+"_"+time.Now().UTC().Format("2006-01-02_15:04:05"))
		fmt.Println("备份命令", backShell)
		retMkdir := exec.Command("bash", "-c", "mkdir -p "+savedir)
		retMkdirBytes, err := retMkdir.Output()
		if (err != nil) {
			fmt.Println("创建目录 出现错误", string(retMkdirBytes), err.Error())
		}

		retFrp := exec.Command("bash", "-c", backShell)
		retFrpBytes, err := retFrp.Output()

		if (err != nil) {
			fmt.Println("出现错误", string(retFrpBytes), err.Error())
		}

		fmt.Println("数据库 ", db ," 备份完毕")
	}
}
