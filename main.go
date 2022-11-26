package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {

	resp, err := http.Get("https://teach-in.ru/course/calculus-shaposhnikov-part1")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()

	tml, _ := io.ReadAll(resp.Body)

	var sb strings.Builder
	sb.Write(tml)

	doc, _ := html.Parse(strings.NewReader(sb.String()))

	videos := []string{}
	titles := []string{}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "meta" && n.Attr[0].Val == "og:video" {
			videos = append(videos, strings.TrimSpace(n.Attr[1].Val))
		}
		if n.Type == html.ElementNode && n.Data == "div" && n.Attr[0].Val == "video-collection-item--info-title" {
			titles = append(titles, strings.TrimSpace(n.FirstChild.Data))
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	for i := 0; i < len(videos); i++ {
		fmt.Println(videos[i])
		resp, _ := http.Get(string(videos[i]))
		out, _ := os.Create(titles[i] + ".mp4")

		io.Copy(out, resp.Body)

		out.Close()
	}
}
