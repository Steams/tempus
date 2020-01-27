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

	// tasks := repo.GetTasksByDay(time.Now())
	// for _, x := range tasks {
	// 	fmt.Println(x.Name)
	// 	fmt.Println(x.Tags)
	// }

	activity_id := repo.NewActivitySession("Work")

	task_session_id := repo.NewTaskSession("Appraisal : Add row border color based on status and ALL filter tab ", activity_id, []string{})

	repo.AddTask(activity.Task{"Appraisal : Add row border color based on status and ALL filter tab ", time.Now().Add(time.Minute * (-40)), time.Now()}, task_session_id)

	return nil
}
