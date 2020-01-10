// +build !dev

package main

import (
	"os"

	"github.com/go-qamel/qamel"
)

func runQtApp(argc int, argv []string) {
	app := qamel.NewApplication(len(os.Args), os.Args)

	app.SetApplicationDisplayName("Tempus")

	// qamel.RegisterQmlListModel("Qamel", 1, 0, "ListModel")

	view := qamel.NewViewerWithSource("qrc:/res/main.qml")
	view.SetResizeMode(qamel.SizeRootObjectToView)
	view.SetHeight(600)
	view.SetWidth(800)
	view.Show()

	app.Exec()
}
