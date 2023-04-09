4chan-rss: main.go
	go mod tidy
	go build -v -ldflags="-s -w"
	ls -lh $@

.PHONY: update
update:
	go list -m -u all
	go get -u
	go mod tidy

.PHONY: clean
clean: ; go clean -x ./...

.PHONY: install
install: 4chan-rss
	upx 4chan-rss
	mv 4chan-rss $(HOME)/.newsboat/feeds/
