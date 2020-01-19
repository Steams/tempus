package activity

import "time"

type Service interface {
	NewTimer(act_name, task_name string) Timer
	GetTasks() []TaskSession
}

type Repository interface {
	AddTask(task Task, session_id string)
	NewActivitySession(name string) string
	NewTaskSession(name, activity_id string) string
	GetTasks() []TaskSession
}

type Timer interface {
	Pause()
	Continue()
	NewTask(string)
	UpdateDuration()
	GetDuration() time.Duration
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
	activity                ActivitySession
	current_task_session_id string
	current_task            task_starter
	duration                time.Duration
	paused                  bool
	repo                    Repository
}

// rename this new activity, each activity gets its own timer instance
func (s service_imp) NewTimer(activity, task_name string) Timer {
	activity_session_id := s.repo.NewActivitySession(activity)
	task_session_id := s.repo.NewTaskSession(task_name, activity_session_id)

	return &timer_impl{
		activity:                ActivitySession{activity_session_id, activity},
		current_task_session_id: task_session_id,
		current_task:            task_starter{task_name, time.Now()},
		duration:                (0 * time.Second),
		paused:                  false,
		repo:                    s.repo,
	}
}

func (s service_imp) GetTasks() []TaskSession {
	return s.repo.GetTasks()
}

func (t *timer_impl) UpdateDuration() {
	t.duration = t.duration + (1 * time.Second)
}

func (t *timer_impl) Pause() {
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
	task := t.current_task.end()
	t.repo.AddTask(task, t.current_task_session_id)

	t.current_task_session_id = t.repo.NewTaskSession(task_name, t.activity.id)
	t.current_task = task_starter{task_name, time.Now()}
}
