> 一键配置Windows vmware开机启动的程序
> 支持windows7
> windows10
> windows server 2012
> windows server 2016

编译运行步骤命令:
```shell
export CGO_ENABLED=1
export CC=x86_64-w64-mingw32-gcc
export CXX=x86_64-w64-mingw32-g++
GOOS=windows GOARCH=amd64 go build -ldflags="-H windowsgui -X" -o vmware-win.exe main.go
```