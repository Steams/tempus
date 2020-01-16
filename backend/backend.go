package backend

import (
	"fmt"
	"os"
	"tempus/pkg/activity"
	activity_repo "tempus/pkg/activity/sqlite_repo"

	_ "github.com/mattn/go-sqlite3"

	"time"

	"github.com/jmoiron/sqlx"

	"github.com/go-qamel/qamel"
)

var service activity.Service

func init() {
	os.Remove("/home/steams/Development/tempus/tempus.db")

	db, err := sqlx.Open("sqlite3", "/home/steams/Development/tempus/tempus.db")

	if err != nil {
		panic(err)
	}

	db.MustExec(activity_repo.Schema)
	repo := activity_repo.New(db)

	service := activity.CreateService(repo)

}

type Backend struct {
	qamel.QmlObject
	_ func(string)         `signal:"timeChanged"`
	_ func()               `signal:"signalPause"`
	_ func()               `signal:"signalStop"`
	_ func()               `signal:"signalStart"`
	_ func()               `signal:"updateList"`
	_ func(string, string) `slot:"toggleStart"`
	_ func(string)         `slot:"changeActivity"`
	_ func(string)         `slot:"load"`

	is_running bool
	is_paused  bool
	timer      activity.Timer
	stopper    chan int
}

func (b *Backend) load() {
	t := activity.NewTimer("Working", "Reading the docs", repo)
	time.Sleep(time.Second * 3)
	t.Pause()
	time.Sleep(time.Second * 3)
	t.Continue()
	t.NewTask("Building out load balancer")
	time.Sleep(time.Second * 3)
	t.Pause()

	stuff := activity.GetTasks()
	for x := range stuff {
		b.updateList(x.activity, x.name)
	}
}

func (b *Backend) changeActivity(name string) {
	b.pause()
	b.is_paused = false
	b.signalStop()
}

func (b *Backend) toggleStart(act_name, task_name string) {
	if b.is_paused {
		b.cont()
		b.signalStart()
		return
	}

	if b.is_running {
		b.pause()
		b.signalPause()
		b.is_paused = true
		return
	}

	b.start(act_name, task_name)
	b.signalStart()

}

func (b *Backend) start(activity_name, task string) {

	b.timer = activity.NewTimer(activity_name, task, repo)
	b.is_running = true
	b.stopper = make(chan int)

	fmt.Println(b.timer)

	go b.runTimer()
}

func (b *Backend) pause() {
	fmt.Println("Pausing")

	b.timer.Pause()
	b.stopper <- 1
	b.is_running = false
}

func (b *Backend) cont() {
	go b.runTimer()
	b.is_paused = false
	b.timer.Continue()
}

func (b *Backend) runTimer() {
	for {
		select {
		case <-b.stopper:
			fmt.Println("Stopping")
			return
		default:
			time.Sleep(time.Second)
			b.timer.UpdateDuration()
			fmt.Println(b.timer)
			t := b.timer.GetDuration()
			b.timeChanged(fmt.Sprintf("%d hrs %d min %d s", int(t.Hours()), int(t.Minutes()), int(t.Seconds())))

		}
	}
}
