# 注意
- 当前只编译了linux版本，如需windows版本可以给一个issue
- 需要安装mysql-client 或者mariadb-client
- 如果是windows平台记得给mysqldump增加一下环境变量

ubuntu/debian系统
```shell
apt install mysql-client
```
centos系统
```shell
yum install mysql-client
```

### 目录结构
- config 所有的配置文件
    - huix.conf.demo 配置文件小样
- main.go 实现的代码文件

### 启动程序
```
$ chmod +x ./backmysql
$ nohup /usr/local/backmysql/./backmysql > /var/log/backmysql.log 2>&1 &
```

### 编译
##### windows开发机
- `set GOOS=linux/windows`
- `set GOARCH=amd64/arm`
- `go build -o backmysql main.go`

##### mac开发机
- `CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go && mv main backmysql`
- `CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -a backmysql main.go`

##### linux开发机
- `CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build main.go`
- `CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build main.go`
## 注意
修改完配置记得`重新启动程序`!!


## 附上安装指定版本的postgresql-client
```bash
sudo sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list'
wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -
sudo apt-get update
# 比如这里就是安装的13版本
apt install postgresql-client-13
```