
.PHONY: linux
linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/bot main.go

.PHONY: mac
mac:
	go build -o ./build/bot_mac main.go

# 编译到 windows
.PHONY: build-windows
windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./build/seek_windows.exe main.go

.PHONY: clean
clean:
	rm -rf ./build

# 编译到 全部平台
.PHONY: build-all
all:
	make clean
	make linux
	make mac




