package main

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"GalgameBox/database"
	"GalgameBox/metadata"
	"GalgameBox/model"
	"GalgameBox/service"

	"github.com/wailsapp/wails/v2/pkg/runtime"
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

// SelectExecutable 打开系统文件选择窗口并返回用户选择的 EXE。
func (a *App) SelectExecutable() (string, error) {
	if a == nil || a.ctx == nil {
		return "", errors.New("应用上下文未初始化")
	}

	return runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择游戏可执行文件",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Windows 可执行文件 (*.exe)",
				Pattern:     "*.exe",
			},
		},
	})
}

func (a *App) StartImport(executablePath string) (service.StartImportResult, error) {
	importService, err := a.requireImportService()
	if err != nil {
		return service.StartImportResult{}, err
	}

	return importService.StartImport(executablePath)
}

func (a *App) PrepareImportMetadata(
	draft service.ImportDraft,
	sources []metadata.Source,
) (service.ImportMetadataResult, error) {
	importService, err := a.requireImportService()
	if err != nil {
		return service.ImportMetadataResult{}, err
	}

	return importService.PrepareImportMetadata(draft, sources)
}

func (a *App) SaveImport(request service.SaveImportRequest) (model.Game, error) {
	importService, err := a.requireImportService()
	if err != nil {
		return model.Game{}, err
	}

	return importService.SaveImport(request)
}

func (a *App) requireImportService() (*service.ImportService, error) {
	if a == nil || a.importService == nil {
		return nil, errors.New("导入服务未初始化")
	}

	return a.importService, nil
}
