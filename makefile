build:
	GOOS=linux GOARCH=amd64 go build -gcflags=-trimpath=$(GOPATH) -asmflags=-trimpath=$(GOPATH) -ldflags "-w -s" -o bin/main_linux main.go
	GOOS=windows GOARCH=amd64 go build -gcflags=-trimpath=$(GOPATH) -asmflags=-trimpath=$(GOPATH) -ldflags "-w -s" -o bin/main_win.exe main.go
compress:
	upx --ultra-brute bin/main_linux
	upx --ultra-brute bin/main_win.exe
package:
	tar zcvf out/dist_linux.tgz public/ migrations/ config.yaml -C bin/ main_linux
	tar zcvf out/dist_win.tgz public/ migrations/ config.yaml -C bin/ main_win.exe