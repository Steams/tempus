package main

import (
	"fmt"
	"os"
	"time"

	activity_repo "tempus/pkg/activity/sqlite_repo"

	_ "github.com/mattn/go-sqlite3"

	"github.com/jmoiron/sqlx"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {

	fmt.Println("Initializing db")

	// os.Remove("./tempus_cli.db")

	db, err := sqlx.Open("sqlite3", "./tempus_cli.db")

	if err != nil {
		panic(err)
	}

	// db.MustExec(activity_repo.Schema)
	repo := activity_repo.New(db)

	// service := activity.CreateService(repo)
	// fmt.Println(service.GetTasks())

	// fmt.Println(repo.GetTasksByDay(time.Now()))
	for _, x := range repo.GetTasksByDay(time.Now().AddDate(0, 0, -1)) {
		fmt.Println("--")
		fmt.Println(x)
		fmt.Println("--")
	}

	return nil
}
