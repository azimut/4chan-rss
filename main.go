package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gorilla/feeds"
	"github.com/k3a/html2text"
	"github.com/moshee/go-4chan-api/api"
	"mvdan.cc/xurls/v2"
)

var URLRegex = xurls.Strict()

var options struct {
	boardName string
	pages     uint
	replies   uint
}

func init() {
	flag.UintVar(&options.replies, "n", 10, "cutoff of number of replies on thread")
	flag.UintVar(&options.pages, "p", 1, "number of pages/request to get/make")
	flag.StringVar(&options.boardName, "b", "g", "board name")
}

func main() {
	if rss, err := run(); err != nil {
		flag.Usage()
		log.Fatal(err)
	} else {
		fmt.Println(rss)
	}
}

func run() (string, error) {
	flag.Parse()
	now := time.Now()
	feed := &feeds.Feed{
		Title: fmt.Sprintf("4chan /%s/ threads", options.boardName),
		Link: &feeds.Link{
			Href: fmt.Sprintf("https://boards.4channel.org/%s/", options.boardName),
		},
		Description: fmt.Sprintf(
			"threads from /%s/ with more than %d comments",
			options.boardName,
			options.replies,
		),
		Author:  &feeds.Author{Name: "Anon"},
		Created: now,
	}
	threads, err := getThreads(options.boardName, options.pages)
	if err != nil {
		return "", err
	}
	feed.Items = processThreads(threads)
	atom, err := feed.ToAtom()
	if err != nil {
		return "", err
	}
	return atom, nil
}

func getThreads(board string, pages uint) (threads []*api.Thread, err error) {
	for i := 0; i < int(pages); i++ {
		newthreads, err := api.GetIndex(board, i)
		if err != nil {
			return nil, err
		}
		threads = append(threads, newthreads...)
	}
	return
}

func processThreads(threads []*api.Thread) []*feeds.Item {
	var items []*feeds.Item
	for _, thread := range threads {
		if thread.Replies() < int(options.replies) {
			continue
		}
		item := processPost(thread.OP)
		item.Title = fmt.Sprintf("[%3d] %s", min(999, thread.Replies()), item.Title)
		items = append(items, item)
	}
	return items
}

func processPost(post *api.Post) *feeds.Item {
	item := &feeds.Item{}
	item.Title = getTitle(post)
	item.Link = &feeds.Link{
		Href: fmt.Sprintf("https://boards.4channel.org/%s/thread/%d/", options.boardName, post.Id),
	}
	item.Description = anchorize(strings.ReplaceAll(post.Comment, "<wbr>", ""))
	if post.File != nil {
		item.Description += fmt.Sprintf("<img alt='%s' src='%s'/>", post.File.Name, post.ImageURL())
	}
	item.Author = &feeds.Author{Name: post.Name}
	item.Created = post.Time
	return item
}

func anchorize(comment string) string {
	return URLRegex.ReplaceAllString(comment, "<a href='$0'>$0</a>")
}

func getTitle(post *api.Post) string {
	title := post.Subject
	if title == "" {
		title = html2text.HTML2Text(post.Comment)
		title = substring(title, 80)
		title = strings.TrimSpace(title)
	}
	if title == "" && post.File != nil {
		title = substring(post.File.Name, 80)
	}
	if title == "" {
		title = "no title"
	}
	return title
}

func substring(s string, end int) string {
	unline := strings.ReplaceAll(s, "\n", " ")
	return unline[:min(len(s), end-1)]
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}
