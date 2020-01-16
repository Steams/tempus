package activity

import "time"

type Service interface {
	NewTimer(act_name, task_name string) Timer
	GetTasks() []DomainTask
}

type Repository interface {
	AddTask(task Task, session_id string)
	NewActivitySession(name string) string
	NewTaskSession(name, activity_id string) string
	GetTasks() []DomainTask
}

type Timer interface {
	Pause()
	Continue()
	NewTask(string)
	UpdateDuration()
	GetDuration() time.Duration
}

type DomainTask struct {
	Id       string
	Name     string
	Act_name string
	// Start    string
	// End      string
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
	id   string
	name string
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
	activity             ActivitySession
	task_sessions        []TaskSession
	current_task_session TaskSession
	current_task         task_starter
	duration             time.Duration
	paused               bool
	repo                 Repository
}

// rename this new activity, each activity gets its own timer instance
func (s service_imp) NewTimer(activity, task_name string) Timer {
	activity_session_id := s.repo.NewActivitySession(activity)
	task_session_id := s.repo.NewTaskSession(task_name, activity_session_id)

	return &timer_impl{
		activity:             ActivitySession{activity_session_id, activity},
		task_sessions:        []TaskSession{},
		current_task_session: TaskSession{task_session_id, task_name},
		current_task:         task_starter{task_name, time.Now()},
		duration:             (0 * time.Second),
		paused:               false,
		repo:                 s.repo,
	}
}

func (s service_imp) GetTasks() []DomainTask {
	return s.repo.GetTasks()
}

func (t *timer_impl) UpdateDuration() {
	t.duration = t.duration + (1 * time.Second)
}

func (t *timer_impl) Pause() {
	task := t.current_task.end()
	// t.tasks = append(t.tasks, t.current_task.end())
	t.repo.AddTask(task, t.current_task_session.id)
	// t.repo.AddTask(t.current_task.end(), t.activity)
}

func (t *timer_impl) Continue() {

	t.current_task = task_starter{t.current_task.name, time.Now()}
}

func (t *timer_impl) GetDuration() time.Duration {
	return t.duration
}

func (t *timer_impl) NewTask(task_name string) {
	// t.tasks = append(t.tasks, t.current_task.end())

	task := t.current_task.end()
	t.repo.AddTask(task, t.current_task_session.id)

	t.current_task_session = TaskSession{t.repo.NewTaskSession(task_name, t.activity.id), task_name}
	t.current_task = task_starter{task_name, time.Now()}
}
