package server

import "html/template"

// Application 結構用來存放共享的依賴，例如模板快取和資料庫連線
type Application struct {
	TemplateCache map[string]*template.Template
}
