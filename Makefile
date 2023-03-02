4chan-rss: main.go
	go build -v -ldflags="-s -w"
	ls -lh 4chan-rss

.PHONY: cache
clean: ; $(GO) clean -x ./...

.PHONY: install
install: 4chan-rss
	install 4chan-rss $(HOME)/go/bin/4chan-rss
