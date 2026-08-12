package main

import (
	"GalgameBox/database"
	"GalgameBox/service"
	"context"
	"database/sql"
	"log"
)

// App struct
type App struct {
	ctx context.Context
	db  *sql.DB

	importService *service.ImportService
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	db, err := database.Open()
	if err != nil {
		log.Printf("打开数据库失败: %v", err)
		return
	}

	a.db = db
	a.importService = service.NewImportService(db)
}

func (a *App) shutdown(ctx context.Context) {
	if a.db != nil {
		if err := a.db.Close(); err != nil {
			log.Printf("关闭数据库失败: %v", err)
		}
	}
}

func (a *App) StartImport(executablePath string) (service.StartImportResult, error) {
	return a.importService.StartImport(executablePath)
}
