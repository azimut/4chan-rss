4chan-rss: main.go
	go mod tidy
	go build -v -ldflags="-s -w"
	ls -lh $@

.PHONY: clean
clean: ; $(GO) clean -x ./...

.PHONY: install
install: 4chan-rss
	upx 4chan-rss
	mv 4chan-rss $(HOME)/.newsboat/feeds/
