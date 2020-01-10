package main

import (
	"os"

	_ "tempus/backend"
)

func main() {
	runQtApp(len(os.Args), os.Args)
}
