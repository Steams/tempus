package sqlite_repo

import (
	"fmt"
	"log"
	"tempus/pkg/activity"

	"github.com/google/uuid"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Repository = activity.Repository
type Task = activity.Task
type DomainTask = activity.DomainTask

const Schema = `
CREATE TABLE activity_session (
    id text,
    name text
);
CREATE TABLE task_session (
    id text,
    name text,
    activity_session_id text
);
CREATE TABLE task (
    name text,
    activity_name text,
    start datetime,
    end datetime,
    duration integer,
    task_session_id
);
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

func (r repository) AddTask(task Task, session_id string) {

	stmt, err := r.db.Preparex("INSERT INTO task(name, start, end, task_session_id) values(?,?,?,?)")
	if err != nil {
		panic(err)
	}
	stmt.MustExec(task.Name, task.Start.String(), task.End.String(), session_id)
}

func (r repository) NewActivitySession(activity string) string {

	stmt, err := r.db.Preparex("INSERT INTO activity_session(id, name) values(?,?)")
	if err != nil {
		panic(err)
	}
	id := uuid.New().String()
	stmt.MustExec(id, activity)
	return id
}

func (r repository) NewTaskSession(name, activity_id string) string {

	stmt, err := r.db.Preparex("INSERT INTO task_session(id, name,activity_session_id) values(?,?,?)")
	if err != nil {
		panic(err)
	}
	id := uuid.New().String()
	stmt.MustExec(id, name, activity_id)
	return id
}

func (r repository) GetTasks() []DomainTask {
	// NOTE i dont like that StructScan requires the name of the struct field to match the name of the response label in the sql
	rows, err := r.db.Queryx("SELECT t.id,t.name,a.name AS act_name FROM task_session AS t INNER JOIN activity_session AS a ON t.activity_session_id = a.id")

	if err != nil {
		log.Fatalln(err)
	}

	m := DomainTask{}
	var tasks []DomainTask

	for rows.Next() {
		err := rows.StructScan(&m)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("%+v\n", m)

		tasks = append(tasks, m)
	}

	if tasks == nil {
		tasks = make([]DomainTask, 0)
		return tasks
	}

	return tasks
}
