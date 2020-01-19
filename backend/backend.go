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

	service = activity.CreateService(repo)
}

type Backend struct {
	qamel.QmlObject
	_ func(string)                                 `signal:"timeChanged"`
	_ func()                                       `signal:"signalPause"`
	_ func()                                       `signal:"signalStop"`
	_ func()                                       `signal:"signalStart"`
	_ func(string, string, string, string, string) `signal:"updateList"`
	_ func()                                       `signal:"clearList"`
	_ func(string, string)                         `slot:"toggleStart"`
	_ func()                                       `slot:"changeActivity"`
	_ func(string)                                 `slot:"changeTask"`
	_ func()                                       `slot:"load"`

	is_running bool
	is_paused  bool
	timer      activity.Timer
	stopper    chan int
}

func (b *Backend) dispatchListUpdate() {
	stuff := service.GetTasks()
	fmt.Println(stuff)
	for _, x := range stuff {
		b.updateList(x.Act_name, x.Name, x.Tasks[0].Start.String(), x.Tasks[0].End.String(), x.Tasks[0].End.Sub(x.Tasks[0].Start).String())
	}
}

func (b *Backend) load() {
	b.dispatchListUpdate()
}

func (b *Backend) changeActivity() {

	if b.is_paused {
		b.is_running = false

		b.dispatchListUpdate()

		b.signalStop()

		return
	}
	if b.is_running {
		b.stop()
	}
}

func (b *Backend) changeTask(name string) {
	b.timer.NewTask(name)

	b.clearList()
	b.dispatchListUpdate()
}

func (b *Backend) toggleStart(act_name, task_name string) {
	if b.is_paused {
		fmt.Println("is paused, now continue")
		b.cont()
		b.signalStart()
		return
	}

	if b.is_running {
		fmt.Println("is running, now pause")
		b.pause()
		b.signalPause()
		return
	}

	fmt.Println("is stop, now start")

	b.start(act_name, task_name)
	b.signalStart()

}

func (b *Backend) start(activity_name, task string) {

	b.timer = service.NewTimer(activity_name, task)
	b.is_running = true
	b.stopper = make(chan int)

	fmt.Println(b.timer)

	go b.runTimer()
}

func (b *Backend) pause() {
	fmt.Println("Pausing")

	b.timer.Pause()
	b.stopper <- 1
	b.is_paused = true
}

func (b *Backend) stop() {
	fmt.Println("Stopping")
	b.timer.Pause()
	b.stopper <- 1
	b.is_paused = false
	b.is_running = false

	b.dispatchListUpdate()
	b.signalStop()
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
