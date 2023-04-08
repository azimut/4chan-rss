4chan-rss: main.go
	go build -v -ldflags="-s -w"
	ls -lh $@

.PHONY: cache
clean: ; $(GO) clean -x ./...

.PHONY: install
install: 4chan-rss
	upx 4chan-rss
	mv 4chan-rss $(HOME)/.newsboat/feeds/
