package crawler_test

import (
	"net/http"
	"net/url"
	"testing"

	"crawler/internal/logger"
	"crawler/internal/mock"

	"github.com/matryer/is"

	"crawler"
)

func TestCrawler_Crawl_noLinksWhenResponseDoesntHaveAnyHRef(t *testing.T) {
	for name, response := range map[string]string{
		"empty page":           ``,
		"json instead of html": `{"name": "Michael"}`,
		"rubbish":              `!#*%!)#u31013`,
		"not a valid html":     `<html></invalid></html>`,
		"empty html page":      `<html></html>`,
		"html page with no links": `<html>
			<head>
				<title>Some test</title>
				<link rel="stylesheet" href="css/theme.css">
			</head>
			<body>
				<h1>Title</h1>
				<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed
				do eiusmod tempor incididunt ut labore et dolore magna aliqua.</p>
			</body>
		</html>`,
	} {
		t.Run(name, func(t *testing.T) {
			is := is.New(t)

			page, err := url.Parse("http://some.url/")
			is.NoErr(err)

			crawler := crawler.Crawler{
				Logger:     logger.Discard(),
				HTTPGetter: mock.NewHTTPGetter(http.StatusOK, response),
			}

			crawler.Crawl(page, func(*url.URL) {
				is.Fail() // this should never be called
			})
		})
	}
}

func TestCrawler_Crawl_returnsEveryLinkOnThePage(t *testing.T) {
	is := is.New(t)

	pageURL := "http://some.url/"

	page, err := url.Parse(pageURL)
	is.NoErr(err)

	want := []string{
		"http://some.url/",
		"http://some.url/ipsum",
		"http://some.url/dolor",
	}

	content := `<html>
		<head>
			<title>Title</title>
		</head>
		<body>
			<h1><a href="/">Home</a></h1>
			<p>Lorem <a href="/ipsum">ipsum</a> <a href="/dolor">dolor</a> sit amet</p>
		</body>
	</html>`

	crawler := crawler.Crawler{
		Logger:     logger.Discard(),
		HTTPGetter: mock.NewHTTPGetter(http.StatusOK, content),
	}

	var got []string
	crawler.Crawl(page, func(url *url.URL) {
		got = append(got, url.String())
	})

	is.Equal(got, want)
}
