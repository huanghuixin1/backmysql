@echo off
set GOOS=windows
set GOARCH=amd64

go build -o ../backmysql.exe ../main.go