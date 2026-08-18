package main

import (
	"testing"

	"GalgameBox/service"
)

func TestImportBridgeReturnsErrorWhenServiceIsNotInitialized(t *testing.T) {
	app := &App{}

	if _, err := app.StartImport("game.exe"); err == nil {
		t.Fatal("StartImport should return an initialization error")
	}
	if _, err := app.PrepareImportMetadata(service.ImportDraft{}, nil); err == nil {
		t.Fatal("PrepareImportMetadata should return an initialization error")
	}
	if _, err := app.SaveImport(service.SaveImportRequest{}); err == nil {
		t.Fatal("SaveImport should return an initialization error")
	}
}
