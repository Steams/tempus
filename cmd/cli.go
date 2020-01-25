package main

import (
	"fmt"
	"os"
	"tempus/pkg/activity"
	activity_repo "tempus/pkg/activity/sqlite_repo"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
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

	// db, err := sqlx.Open("sqlite3", "./tempus_cli.db")

	db, err := sqlx.Open("sqlite3", "./tempus.db")

	if err != nil {
		panic(err)
	}

	// db.MustExec(activity_repo.Schema)
	repo := activity_repo.New(db)

	activity_id := repo.NewActivitySession("Project")

	task_session_id := repo.NewTaskSession("Tempus : refactoring tag id caching system", activity_id, []string{"Go", "Tempus"})

	repo.AddTask(activity.Task{"Tempus : refactoring tag id caching system", time.Now().Add(time.Minute * (-40)), time.Now()}, task_session_id)

	return nil
}
