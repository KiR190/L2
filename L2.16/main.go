package main

import (
	"flag"
	"log"

	"mywget/downloader"
)

func main() {
	url := flag.String("url", "https://example.com", "URL сайта для скачивания")
	output := flag.String("o", "site_copy", "Папка для сохранения сайта")
	depth := flag.Int("depth", 3, "Глубина рекурсии")
	timeout := flag.Int("timeout", 15, "Таймаут HTTP-запросов (сек)")
	parallel := flag.Int("parallel", 5, "Количество параллельных загрузок")
	robots := flag.Bool("robots", true, "Уважать robots.txt")

	flag.Parse()

	if *url == "" {
		log.Fatal("Не указан URL. Используйте -url https://example.com")
	}

	d := downloader.NewDownloader(*url, *output, *robots, *depth, *parallel, *timeout)
	d.Run(*url)

	log.Println("Скачивание завершено.")
}
