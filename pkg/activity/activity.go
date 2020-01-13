package activity

import "time"

type task struct {
	name  string
	start time.Time
	end   time.Time
}

type task_starter struct {
	name  string
	start time.Time
}

func (t task_starter) end() task {
	return task{t.name, t.start, time.Now()}
}

type Timer interface {
	Pause()
	Continue()
	NewTask(string)
	UpdateDuration()
	GetDuration() time.Duration
}

type timer_impl struct {
	activity     string
	tasks        []task
	current_task task_starter
	duration     time.Duration
	paused       bool
}

func NewTimer(activity, task_name string) Timer {
	return &timer_impl{
		activity:     activity,
		tasks:        []task{},
		current_task: task_starter{task_name, time.Now()},
		duration:     (0 * time.Second),
		paused:       false,
	}
}

func (t *timer_impl) UpdateDuration() {
	t.duration = t.duration + (1 * time.Second)
}

func (t *timer_impl) Pause() {
	t.tasks = append(t.tasks, t.current_task.end())
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
