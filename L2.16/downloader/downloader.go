package downloader

import (
	"log"
	"os"
	"strings"
	"time"

	"mywget/parser"
	"mywget/utils"

	"github.com/gocolly/colly/v2"
)

type Downloader struct {
	//Visited       map[string]bool
	OutputDir     string
	Collector     *colly.Collector
	Domain        string
	RespectRobots bool
	MaxDepth      int
	Parallelism   int
}

func NewDownloader(startURL, outputDir string, respectRobots bool, maxDepth int, parallelism int, timeout int) *Downloader {
	domain := utils.GetDomain(startURL)

	opts := []colly.CollectorOption{
		colly.AllowedDomains(domain),
		colly.MaxDepth(maxDepth),
		colly.Async(true),
	}

	if !respectRobots {
		opts = append(opts, colly.IgnoreRobotsTxt())
	}

	c := colly.NewCollector(opts...)

	c.SetRequestTimeout(time.Duration(timeout) * time.Second)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: parallelism,
	})

	return &Downloader{
		//Visited:       make(map[string]bool),
		OutputDir:     outputDir,
		Collector:     c,
		Domain:        domain,
		RespectRobots: respectRobots,
		MaxDepth:      maxDepth,
		Parallelism:   parallelism,
	}
}

func (d *Downloader) Run(startURL string) {
	c := d.Collector

	// Логирование
	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting:", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Error:", r.Request.URL, err)
	})

	// Переход по ссылкам
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		e.Request.Visit(link)
	})

	// Ресурсы
	c.OnHTML("img[src], script[src], link[rel=stylesheet]", func(e *colly.HTMLElement) {
		src := e.Attr("src")
		if src == "" {
			src = e.Attr("href")
		}
		abs := e.Request.AbsoluteURL(src)
		e.Request.Visit(abs)
	})

	// Сохранение не-HTML ресурсов
	c.OnResponse(func(r *colly.Response) {
		contentType := r.Headers.Get("Content-Type")
		if strings.Contains(contentType, "text/html") {
			return
		}

		localPath := utils.URLToFilePath(d.OutputDir, r.Request.URL)
		if err := utils.EnsureDir(localPath); err != nil {
			log.Println(err)
			return
		}

		if err := os.WriteFile(localPath, r.Body, 0644); err != nil {
			log.Println("Failed to save:", r.Request.URL, err)
		} else {
			log.Println("Saved:", localPath)
		}
	})

	// Сохранение HTML после обработки
	c.OnHTML("html", func(e *colly.HTMLElement) {
		localPath := utils.URLToFilePath(d.OutputDir, e.Request.URL)
		html := parser.RewriteLinks(e, d.OutputDir)

		if err := utils.EnsureDir(localPath); err == nil {
			if err := os.WriteFile(localPath, []byte(html), 0644); err != nil {
				log.Println("Failed to save page:", e.Request.URL, err)
			} else {
				log.Println("Saved page:", localPath)
			}
		} else {
			log.Println(err)
		}
	})

	c.Visit(startURL)
	c.Wait()
}
