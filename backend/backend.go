package backend

import (
	"fmt"
	"tempus/pkg/activity"
	activity_repo "tempus/pkg/activity/sqlite_repo"

	_ "github.com/mattn/go-sqlite3"

	"time"

	"github.com/jmoiron/sqlx"

	"github.com/go-qamel/qamel"
)

var repo activity.Repository

func init() {
	db, err := sqlx.Open("sqlite3", "/home/steams/Development/tempus/tempus.db")
	if err != nil {
		panic(err)
	}

	db.MustExec(activity_repo.Schema)
	repo = activity_repo.New(db)
}

type BackEnd struct {
	qamel.QmlObject
	_ func(string)                `signal:"timeChanged"`
	_ func() bool                 `slot:"isRunning"`
	_ func(string, string) string `slot:"toggleStart"`

	is_running bool
	is_paused  bool
	timer      activity.Timer
	stopper    chan int
}

func (b *BackEnd) isRunning() bool {
	return b.is_running
}

func (b *BackEnd) toggleStart(act_name, task_name string) string {
	if b.is_paused {
		b.cont()
		return "Pause Timer"
	}

	if b.is_running {
		b.pause()
		b.is_paused = true
		return "Continue Timer"
	}

	b.start(act_name, task_name)
	return "Pause Timer"

}

func (b *BackEnd) start(activity_name, task string) {

	b.timer = activity.NewTimer(activity_name, task, repo)
	b.is_running = true
	b.stopper = make(chan int)

	fmt.Println(b.timer)

	go b.runTimer()
}

func (b *BackEnd) pause() {
	fmt.Println("Pausing")

	b.timer.Pause()
	b.stopper <- 1
	b.is_running = false
}

func (b *BackEnd) cont() {
	go b.runTimer()
	b.is_paused = false
}

func (b *BackEnd) runTimer() {
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
