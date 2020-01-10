// +build dev

package main

import (
	"log"
	"os"
	fp "path/filepath"

	"github.com/go-qamel/qamel"
)

func runQtApp(argc int, argv []string) {
	log.Println("DEV MODE")

	// Create QT app
	app := qamel.NewApplication(len(os.Args), os.Args)
	app.SetApplicationDisplayName("Tempus")

	// Register qamel model
	// qamel.RegisterQmlListModel("Qamel", 1, 0, "ListModel")
	// RegisterQmlBackEnd("BackEnd", 1, 0, "BackEnd")

	// 	// Create a QML viewer
	view := qamel.NewViewerWithSource("res/main.qml")
	view.SetResizeMode(qamel.SizeRootObjectToView)
	view.SetHeight(600)
	view.SetWidth(800)
	view.Show()

	// Watch change in resource dir
	projectDir, err := os.Getwd()
	if err != nil {
		log.Fatalln("Failed to get working directory:", err)
	}

	resDir := fp.Join(projectDir, "res")
	go view.WatchResourceDir(resDir)

	// Exec app
	app.Exec()
}
