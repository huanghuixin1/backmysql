package backmysql

import (
	"io/ioutil"
	"strings"
	"time"
	"github.com/robfig/config"
)

func main() {
	// 获取程序的配置
	c, _ := config.ReadDefault("./config.conf")

}