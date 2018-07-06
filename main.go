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
	// 获取程序的配置
	//c, _ := config.ReadDefault("./config.conf")
	files, _ := ioutil.ReadDir("./config")
	for _, file := range files {
		if (!strings.HasSuffix(file.Name(), ".conf")) {
			continue
		}

		configSubProcess, _ := config.ReadDefault("./config/" + file.Name())
		startBackInterval(configSubProcess)
	}
}

func startBackInterval(config *config.Config) {
	user, _ := config.String("", "user")
	pwd, _ := config.String("", "pwd")
	host, _ := config.String("", "host")
	port, _ := config.String("", "port")
	savedir, _ := config.String("", "savedir")
	dbsStr, _ := config.String("", "dbs")
	dbs := strings.Split(dbsStr, ",")

	backtime, err := config.Float("", "backtime")
	// 如果是不float类型 则是时间判断
	if err != nil {
		backTimeStr, _ := config.String("", "backtime")
		hourMinuts := strings.Split(backTimeStr, ":")
		housr, _ := strconv.Atoi(hourMinuts[0])
		minuts, _ := strconv.Atoi(hourMinuts[1])
		for true {
			now := time.Now()
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

}

func invokeBack(user string, pwd string, host string, port string, savedir string, dbs []string) {
	for _, db := range dbs {
		//now := time.Now();
		backShell := fmt.Sprintf("mysqldump --host %s --port %s -u%s -p%s --databases %s > %s%s.sql", host, port, user, pwd, db, savedir, db+"_"+time.Now().Format("2006-01-02_15:04:05"))
		fmt.Println("备份命令", backShell)
		retFrp := exec.Command("bash", "-c", backShell)
		retFrpBytes, errFrp := retFrp.Output()

		if (errFrp != nil) {
			fmt.Println("出现错误", string(retFrpBytes), errFrp.Error())
		}
	}

}
