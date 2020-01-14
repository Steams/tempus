package activity

import "time"

type Timer interface {
	Pause()
	Continue()
	NewTask(string)
	UpdateDuration()
	GetDuration() time.Duration
}

type Repository interface {
	Add(task Task, activity string)
}

type Task struct {
	Name  string
	Start time.Time
	End   time.Time
}

type task_starter struct {
	name  string
	start time.Time
}

func (t task_starter) end() Task {
	return Task{t.name, t.start, time.Now()}
}

type timer_impl struct {
	activity     string
	tasks        []Task
	current_task task_starter
	duration     time.Duration
	paused       bool
	repo         Repository
}

func NewTimer(activity, task_name string, repo Repository) Timer {
	return &timer_impl{
		activity:     activity,
		tasks:        []Task{},
		current_task: task_starter{task_name, time.Now()},
		duration:     (0 * time.Second),
		paused:       false,
		repo:         repo,
	}
}

func (t *timer_impl) UpdateDuration() {
	t.duration = t.duration + (1 * time.Second)
}

func (t *timer_impl) Pause() {
	t.tasks = append(t.tasks, t.current_task.end())
	t.repo.Add(t.current_task.end(), t.activity)
}

func (t *timer_impl) Continue() {
}

func (t *timer_impl) GetDuration() time.Duration {
	return t.duration
}

func (t *timer_impl) NewTask(task_name string) {
	t.tasks = append(t.tasks, t.current_task.end())
	t.current_task = task_starter{task_name, time.Now()}
}
