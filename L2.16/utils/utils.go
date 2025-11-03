package utils

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// Получаем домен из URL
func GetDomain(rawURL string) string {
	u, _ := url.Parse(rawURL)
	return u.Host
}

// Создаем локальный путь для сохранения файла
func URLToFilePath(baseDir string, u *url.URL) string {
	path := u.Path
	if path == "" || strings.HasSuffix(path, "/") {
		path += "index.html"
	}
	path = strings.TrimPrefix(path, "/")
	return filepath.Join(baseDir, path)
}

// Создаем директорию, если не существует
func EnsureDir(filePath string) error {
	dir := filepath.Dir(filePath)
	return os.MkdirAll(dir, 0755)
}

// Парсим URL с обработкой ошибок
func ParseURL(raw string) *url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		return &url.URL{}
	}
	return u
}

// Конвертируем внешний URL в локальный путь для ссылки в HTML
func URLToLocalPath(baseDir string, u *url.URL) string {
	return filepath.ToSlash(strings.TrimPrefix(URLToFilePath(baseDir, u), baseDir+"/"))
}
