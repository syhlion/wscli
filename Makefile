WSCLI:= wscli
TAG := `git describe --tags `.`git rev-parse --short HEAD`
TZ := Asia/Taipei
DATETIME := `TZ=$(TZ) date +%Y/%m/%d.%T`
show-tag:
	echo $(TAG)
verify-glide:
	if [ ! -e `which glide` ] ; then\
		echo 'please install "https://github.com/Masterminds/glide"';\
		exit 1;\
	fi
buildwscli = GOOS=$(1) GOARCH=$(2) go build -ldflags "-X main.version=$(TAG) -X main.name=$(WSCLI) -X main.compileDate=$(DATETIME)($(TZ))" -a -o build/$(WSCLI)$(3) 
tar = cp README.md build/ &&cd build && tar -zcvf $(WSCLI)_$(TAG)_$(1)_$(2).tar.gz $(WSCLI)$(3) README.md && rm $(WSCLI)$(3) README.md

build/linux: 
	go test
	$(call buildwscli,linux,amd64,)
build/linux_amd64.tar.gz: build/linux
	$(call tar,linux,amd64,)
build/windows: 
	go test
	$(call buildwscli,windows,amd64,.exe)
build/windows_amd64.tar.gz: build/windows
	$(call tar,windows,amd64,.exe)
build/darwin: 
	go test
	$(call buildwscli,darwin,amd64,)
build/darwin_amd64.tar.gz: build/darwin
	$(call tar,darwin,amd64,)
clean:
	rm -rf build/
