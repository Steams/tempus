package activity

import (
	"time"
)

type Service interface {
	NewTimer(act_name, task_name string) Timer
	GetTasks() []TaskSession
	GetTasksByDay(time.Time) []TaskSession
}

type Repository interface {
	AddTask(task Task, session_id string)
	NewActivitySession(name string) string
	NewTaskSession(name, activity_id string, tags []string) string
	GetTasks() []TaskSession
	GetActivities() []Activity
	GetTasksByDay(date time.Time) []TaskSession
	GetTags() []Tag
}

type Timer interface {
	Pause()
	Continue()
	NewTask(string)
	UpdateDuration()
	GetDuration() time.Duration
	SetTags([]string)
}

type Tag struct {
	Id   string
	Name string
}

type Task struct {
	Name  string
	Start time.Time
	End   time.Time
}

type ActivitySession struct {
	id   string
	name string
}

type Activity struct {
	Id    string
	Name  string
	Tasks []TaskSession
}

type TaskSession struct {
	Id       string
	Name     string
	Act_name string
	Tasks    []Task
}

type task_starter struct {
	name  string
	start time.Time
}

func (t task_starter) end() Task {
	return Task{t.name, t.start, time.Now()}
}

type service_imp struct {
	repo Repository
}

func CreateService(r Repository) Service {
	return service_imp{r}
}

type timer_impl struct {
	activity_id             string
	activity_name           string
	current_task_session_id string
	current_task            task_starter
	duration                time.Duration
	paused                  bool
	repo                    Repository
	current_tags            []string
}

// rename this new activity, each activity gets its own timer instance
func (s service_imp) NewTimer(activity, task_name string) Timer {
	return &timer_impl{
		activity_id:             "",
		activity_name:           activity,
		current_task_session_id: "",
		current_task:            task_starter{task_name, time.Now()},
		duration:                (0 * time.Second),
		paused:                  false,
		repo:                    s.repo,
		current_tags:            []string{},
	}
}

func (s service_imp) GetTasks() []TaskSession {
	return s.repo.GetTasks()
}

func (s service_imp) GetTasksByDay(date time.Time) []TaskSession {
	return s.repo.GetTasksByDay(date)
}

func (s service_imp) GetActivities() []Activity {
	return s.repo.GetActivities()
}

func (t *timer_impl) UpdateDuration() {
	t.duration = t.duration + (1 * time.Second)
}

func (t *timer_impl) SetTags(tags []string) {
	t.current_tags = tags
}

func (t *timer_impl) Pause() {

	if t.activity_id == "" {
		t.activity_id = t.repo.NewActivitySession(t.activity_name)
	}

	if t.current_task_session_id == "" {
		t.current_task_session_id = t.repo.NewTaskSession(t.current_task.name, t.activity_id, t.current_tags)
	}

	task := t.current_task.end()
	t.repo.AddTask(task, t.current_task_session_id)
}

func (t *timer_impl) Continue() {
	t.current_task = task_starter{t.current_task.name, time.Now()}
}

func (t *timer_impl) GetDuration() time.Duration {
	return t.duration
}

func (t *timer_impl) NewTask(task_name string) {

	if t.activity_id == "" {
		t.activity_id = t.repo.NewActivitySession(t.activity_name)
	}

	if t.current_task_session_id == "" {
		t.current_task_session_id = t.repo.NewTaskSession(t.current_task.name, t.activity_id, t.current_tags)
	}

	task := t.current_task.end()
	t.repo.AddTask(task, t.current_task_session_id)

	t.current_task_session_id = t.repo.NewTaskSession(task_name, t.activity_id, t.current_tags)
	t.current_task = task_starter{task_name, time.Now()}
}
