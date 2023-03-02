package crawler

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"crawler/internal/logger"
	"crawler/internal/workerpool"

	"golang.org/x/net/html"
)

type Crawler struct {
	Logger     logger.Logger
	HTTPGetter interface {
		Get(url string) (*http.Response, error)
	}
}

func (c Crawler) Crawl(startURL *url.URL, found func(*url.URL)) {
	var visited sync.Map

	c.crawl(startURL.Host, startURL, found, &visited)
}

func (c Crawler) crawl(host string, page *url.URL, found func(*url.URL), visited *sync.Map) {
	c.Logger.Infof("crawling %s", page)

	links := c.extractHostLinks(host, page)

	pool := workerpool.New(5)
	defer pool.Stop()

	for _, link := range links {
		if _, ok := visited.LoadOrStore(link, nil); ok {
			// skip already visited links
			continue
		}

		url, err := url.Parse(link)
		if err != nil {
			c.Logger.Errorf("invalid link %s: %v", link, err)
			continue
		}

		found(url)

		pool.Run(func(context.Context) {
			c.crawl(host, url, found, visited)
		})
	}

	pool.Wait()
}

func (c Crawler) extractHostLinks(host string, page *url.URL) []string {
	res, err := c.HTTPGetter.Get(page.String())
	if err != nil {
		c.Logger.Errorf(`can't read url %q: %v`, page, err)
		return nil
	}
	defer res.Body.Close()

	doc, err := html.Parse(res.Body)
	if err != nil {
		c.Logger.Errorf(`can't parse %q content as html: %v`, page, err)
		return nil
	}

	var links []string

	forEachHTMLNode(doc, func(node *html.Node) {
		href := extractHRef(node)
		if href == "" {
			return
		}

		link, err := page.Parse(href)
		if err != nil {
			c.Logger.Errorf("can't parse link %q as an url: %v", href, err)
			return
		}

		if link.Host != host {
			return
		}

		links = append(links, link.String())
	})

	return links
}

func extractHRef(node *html.Node) string {
	if node.Type != html.ElementNode || node.Data != "a" {
		return ""
	}
	for _, a := range node.Attr {
		if a.Key == "href" {
			return strings.ToLower(strings.TrimSpace(a.Val))
		}
	}
	return ""
}

func forEachHTMLNode(node *html.Node, fn func(node *html.Node)) {
	fn(node)
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		forEachHTMLNode(child, fn)
	}
}
