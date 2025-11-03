package parser

import (
	"bytes"
	"path/filepath"
	"strings"

	"mywget/utils"

	"github.com/gocolly/colly/v2"
	"golang.org/x/net/html"
)

// Получаем ресурсы со страницы: img, css, js
func CollectResources(e *colly.HTMLElement) []string {
	var resources []string

	e.ForEach("img[src]", func(_ int, el *colly.HTMLElement) {
		resources = append(resources, el.Request.AbsoluteURL(el.Attr("src")))
	})

	e.ForEach("link[rel='stylesheet']", func(_ int, el *colly.HTMLElement) {
		resources = append(resources, el.Request.AbsoluteURL(el.Attr("href")))
	})

	e.ForEach("script[src]", func(_ int, el *colly.HTMLElement) {
		resources = append(resources, el.Request.AbsoluteURL(el.Attr("src")))
	})

	return resources
}

// Собираем ссылки на другие страницы
func CollectLinks(e *colly.HTMLElement, allowedDomain string) []string {
	var links []string
	e.ForEach("a[href]", func(_ int, el *colly.HTMLElement) {
		link := el.Request.AbsoluteURL(el.Attr("href"))
		if strings.Contains(link, allowedDomain) {
			links = append(links, link)
		}
	})
	return links
}

func RewriteLinks(e *colly.HTMLElement, baseDir string) string {
	doc, err := html.Parse(bytes.NewReader(e.Response.Body))
	if err != nil {
		return string(e.Response.Body)
	}

	currentPath := utils.URLToFilePath(baseDir, e.Request.URL)

	makeRelative := func(targetURL string) string {
		u := utils.ParseURL(e.Request.AbsoluteURL(targetURL))
		targetLocal := utils.URLToFilePath(baseDir, u)
		rel, err := filepath.Rel(filepath.Dir(currentPath), targetLocal)
		if err != nil {
			return targetLocal
		}
		return filepath.ToSlash(rel)
	}

	var rewrite func(*html.Node)
	rewrite = func(n *html.Node) {
		if n.Type == html.ElementNode {
			for i, attr := range n.Attr {
				switch strings.ToLower(attr.Key) {
				case "src", "href":
					abs := e.Request.AbsoluteURL(attr.Val)
					u := utils.ParseURL(abs)
					if n.Data == "a" && u.Host != e.Request.URL.Host {
						continue
					}
					local := makeRelative(attr.Val)
					n.Attr[i].Val = local
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			rewrite(c)
		}
	}

	rewrite(doc)

	var buf bytes.Buffer
	if err := html.Render(&buf, doc); err != nil {
		return string(e.Response.Body)
	}
	return buf.String()
}
