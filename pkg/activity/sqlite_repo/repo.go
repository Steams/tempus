package sqlite_repo

import (
	"fmt"
	"log"
	"tempus/pkg/activity"
	"time"

	"github.com/google/uuid"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Repository = activity.Repository
type Task = activity.Task
type TaskSession = activity.TaskSession
type ActivitySession = activity.Activity

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
    start integer,
    end integer,
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
	stmt.MustExec(task.Name, task.Start.Unix(), task.End.Unix(), task.End.Sub(task.Start), session_id)
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

type TaskScanner struct {
	Name  string
	Start int64
	End   int64
}

func (r repository) GetTasks() []TaskSession {
	session_results, err := r.db.Queryx("SELECT t.id,t.name,a.name AS act_name FROM task_session AS t INNER JOIN activity_session AS a ON t.activity_session_id = a.id")

	if err != nil {
		log.Fatalln(err)
	}

	return extractTaskSessions(session_results, r.db)
}

func (r repository) GetTasksByActivity(activity_session_id string) []TaskSession {
	session_results, err := r.db.Queryx("SELECT t.id,t.name,a.name AS act_name FROM task_session AS t INNER JOIN activity_session AS a ON t.activity_session_id = a.id WHERE a.id = $1", activity_session_id)

	if err != nil {
		log.Fatalln(err)
	}

	return extractTaskSessions(session_results, r.db)
}

func (r repository) GetTasksByDay(date time.Time) []TaskSession {
	year, month, day := date.Date()
	day_start := time.Date(year, month, day, 0, 0, 0, 0, date.Location())
	day_end := time.Date(year, month, day, 24, 0, 0, 0, date.Location())

	fmt.Println(day_start.Unix())
	fmt.Println(day_end.Unix())

	session_results, err := r.db.Queryx(`
	SELECT t.id,t.name,a.name AS act_name
	FROM task_session AS t
	INNER JOIN activity_session AS a
	  ON t.activity_session_id = a.id
	INNER JOIN task
	  ON task.task_session_id = t.id
	WHERE task.start BETWEEN $1 AND $2
	GROUP BY t.id
	`, day_start.Unix(), day_end.Unix())

	if err != nil {
		log.Fatalln(err)
	}

	return extractTaskSessions(session_results, r.db)
}

func extractTaskSessions(results *sqlx.Rows, db *sqlx.DB) []TaskSession {

	sessions := []TaskSession{}

	for results.Next() {
		session := TaskSessionScanner{}

		// NOTE i dont like that StructScan requires the name of the struct field to match the name of the response label in the sql
		err := results.StructScan(&session)

		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("%+v\n", session)

		tasks := []Task{}

		task_results, err := db.Queryx("SELECT name,start,end FROM task WHERE task_session_id = $1", session.Id)

		for task_results.Next() {
			task := TaskScanner{}

			err := task_results.StructScan(&task)

			if err != nil {
				log.Fatalln(err)
			}
			fmt.Printf("%+v\n", task)
			tasks = append(tasks, Task{task.Name, time.Unix(task.Start, 0), time.Unix(task.End, 0)})
		}

		sessions = append(sessions, TaskSession{session.Id, session.Name, session.Act_name, tasks})
	}

	if sessions == nil {
		sessions = make([]TaskSession, 0)
		return sessions
	}

	return sessions
}

func (r repository) GetActivities() []ActivitySession {
	session_results, err := r.db.Queryx("SELECT id, name FROM activity_session")

	if err != nil {
		log.Fatalln(err)
	}

	sessions := []ActivitySession{}

	for session_results.Next() {

		session := struct {
			id   string
			name string
		}{}

		err := session_results.StructScan(&session)

		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("%+v\n", session)

		tasks := r.GetTasksByActivity(session.id)

		sessions = append(sessions, ActivitySession{session.id, session.name, tasks})
	}

	if sessions == nil {
		sessions = make([]ActivitySession, 0)
		return sessions
	}

	return sessions
}
