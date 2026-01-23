package view

import (
	"html/template"
	"path/filepath"
)

// NewCache 解析 templates 目錄下所有 .html 檔案並回傳一個快取 map
func NewCache() (map[string]*template.Template, error) {
	// 初始化一個新的 map 來作為快取
	cache := map[string]*template.Template{}

	// 使用 Glob 找出所有符合 *.html 模式的檔案路徑
	pages, err := filepath.Glob("templates/*.html")
	if err != nil {
		return nil, err
	}

	// 遍歷所有找到的頁面
	for _, page := range pages {
		// 取得檔案名稱作為 key，例如 "index.html"
		name := filepath.Base(page)

		// 解析模板檔案
		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// 將解析好的模板存入快取
		cache[name] = ts
	}

	// 回傳快取
	return cache, nil
}
