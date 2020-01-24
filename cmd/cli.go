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

	// service := activity.CreateService(repo)
	// fmt.Println(service.GetTasks())

	// fmt.Println(uuid.New().String())
	time_layout := "03:04PM"

	// t1, _ := time.Parse(time_layout, "12:00PM")
	// t2, _ := time.Parse(time_layout, "01:40PM")

	year, month, day := time.Now().Date()
	day_start := time.Date(year, month, day, 0, 0, 0, 0, time.Now().Location())
	// morning := day_start.Add(time.Hour * 8)
	start := day_start.Add(time.Hour * 12)
	end := start.Add(time.Hour * 1).Add(time.Minute * 40)

	task := activity.Task{"Appraisal: debugging filter issue", start, end}

	repo.AddTask(task, "fcba795b-5c4a-4092-9435-59909997f748")

	// fmt.Println(morning.Format(time_layout))
	fmt.Println(start.Format(time_layout))
	fmt.Println(end.Format(time_layout))
	// fmt.Println(time.Now().Format(time_layout))
	// fmt.Println(t1.Format(time_layout))
	// fmt.Println(t2.Format(time_layout))

	// fmt.Println(morning.Unix())
	fmt.Println(start.Unix())
	fmt.Println(end.Unix())
	// fmt.Println(time.Now().Unix())
	// fmt.Println(t1.Unix())
	// fmt.Println(t2.Unix())
	// fmt.Println(repo.GetTasksByDay(time.Now()))
	// for _, x := range repo.GetTasksByDay(time.Now().AddDate(0, 0, -1)) {
	// 	fmt.Println("--")
	// 	fmt.Println(x)
	// 	fmt.Println("--")
	// }

	return nil
}
