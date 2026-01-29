package server

import (
	"html/template"
	"web-demo/config"

	"gorm.io/gorm"
)

// Application 結構用來存放共享的依賴
type Application struct {
	Config        *config.AppConfig
	DB            *gorm.DB
	TemplateCache map[string]*template.Template
}
