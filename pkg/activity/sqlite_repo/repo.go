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
type TaskSession = activity.TaskSession

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

func (r repository) AddTask(task Task, session_id string) {

	stmt, err := r.db.Preparex("INSERT INTO task(name, start, end, duration, task_session_id) values(?,?,?,?,?)")
	if err != nil {
		panic(err)
	}
	stmt.MustExec(task.Name, task.Start.String(), task.End.String(), task.End.Sub(task.Start), session_id)
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

type TaskSessionScanner struct {
	Id       string
	Name     string
	Act_name string
}

func (r repository) GetTasks() []TaskSession {
	// NOTE i dont like that StructScan requires the name of the struct field to match the name of the response label in the sql
	session_results, err := r.db.Queryx("SELECT t.id,t.name,a.name AS act_name FROM task_session AS t INNER JOIN activity_session AS a ON t.activity_session_id = a.id")

	if err != nil {
		log.Fatalln(err)
	}

	sessions := []TaskSession{}

	for session_results.Next() {
		session := TaskSessionScanner{}

		err := session_results.StructScan(&session)

		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("%+v\n", session)

		tasks := []Task{}

		task_results, err := r.db.Queryx("SELECT name,start,end FROM task WHERE task_session_id = $1", session.Id)

		for task_results.Next() {
			task := Task{}

			err := task_results.StructScan(&task)

			if err != nil {
				log.Fatalln(err)
			}
			fmt.Printf("%+v\n", task)
			tasks = append(tasks, task)
		}

		sessions = append(sessions, TaskSession{session.Id, session.Name, session.Act_name, tasks})
	}

	if sessions == nil {
		sessions = make([]TaskSession, 0)
		return sessions
	}

	return sessions
}
