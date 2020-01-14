package sqlite_repo

import (
	"tempus/pkg/activity"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Repository = activity.Repository
type Task = activity.Task

const Schema = `
CREATE TABLE task (
    name text,
    start text,
    end text,
    activity text
)
`

func New(db *sqlx.DB) Repository {
	return repository{db}
}

type repository struct {
	db *sqlx.DB
}

func (r repository) Add(task Task, activity string) {

	stmt, err := r.db.Preparex("INSERT INTO task(name, start, end, activity) values(?,?,?,?)")
	if err != nil {
		panic(err)
	}
	stmt.MustExec(task.Name, task.Start.String(), task.End.String(), activity)
}
