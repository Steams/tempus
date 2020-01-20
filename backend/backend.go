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

var service activity.Service
var time_layout string

func init() {
	time_layout = "03:04PM"
	// os.Remove("/home/steams/Development/tempus/tempus.db")

	db, err := sqlx.Open("sqlite3", "/home/steams/Development/tempus/tempus.db")

	if err != nil {
		panic(err)
	}

	// db.MustExec(activity_repo.Schema)
	repo := activity_repo.New(db)

	service = activity.CreateService(repo)
}

type Backend struct {
	qamel.QmlObject
	_ func(string)                                   `signal:"timeChanged"`
	_ func()                                         `signal:"signalPause"`
	_ func()                                         `signal:"signalStop"`
	_ func()                                         `signal:"signalStart"`
	_ func(string, string, string, string, string)   `signal:"updateList"`
	_ func(string, string, float64, string, float64) `signal:"updateTimeline"`
	_ func(string, string, float64)                  `signal:"updateReport"`
	_ func()                                         `signal:"clearList"`
	_ func(string, string)                           `slot:"toggleStart"`
	_ func()                                         `slot:"changeActivity"`
	_ func(string)                                   `slot:"changeTask"`
	_ func()                                         `slot:"load"`

	is_running bool
	is_paused  bool
	timer      activity.Timer
	stopper    chan int
}

func duration(tasks []activity.Task) string {
	var sum time.Duration = 0
	for _, x := range tasks {
		sum = sum + x.End.Sub(x.Start)
	}
	return sum.String()
}

func (b *Backend) dispatchListUpdate() {
	stuff := service.GetTasks()
	fmt.Println(stuff)
	for _, x := range stuff {
		b.updateList(
			x.Act_name,
			x.Name,
			x.Tasks[0].Start.Format(time_layout),
			x.Tasks[len(x.Tasks)-1].End.Format(time_layout),
			duration(x.Tasks),
		)
	}
}

func (b *Backend) dispatchReportUpdate() {
	// TODO this needs to be get activities BY TYPE "work ...", or u can get all activities and aggregate them into groups urslef
	// stuff := service.GetActivities()
	// fmt.Println(stuff)
	// for _, x := range stuff {
	// 	b.updateList(
	// 		x.Act_name,
	// 		x.Name,
	// 		x.Tasks[0].Start.Format(time_layout),
	// 		x.Tasks[len(x.Tasks)-1].End.Format(time_layout),
	// 		duration(x.Tasks),
	// 	)
	// }
}

func (b *Backend) dispatchTimelineUpdate() {
	stuff := service.GetTasks()
	fmt.Println("Timeline stuff")
	fmt.Println(stuff)

	year, month, day := time.Now().Date()
	day_start := time.Date(year, month, day, 0, 0, 0, 0, time.Now().Location())
	morning := day_start.Add(time.Hour * 8)

	fmt.Println(morning.Format(time_layout))

	for _, i := range stuff {
		for _, x := range i.Tasks {

			b.updateTimeline(
				x.Start.Format(time_layout),
				x.End.Format(time_layout),
				x.End.Sub(x.Start).Hours(),
				i.Act_name+" : "+x.Name,
				x.Start.Sub(morning).Hours(),
			)
		}
	}

}

func (b *Backend) load() {
	b.dispatchListUpdate()
	// b.dispatchReportUpdate()
	b.dispatchTimelineUpdate()
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
			b.timeChanged(fmt.Sprintf("%d hrs %d min %d s", int(t.Hours()), int(t.Minutes())%60, int(t.Seconds())%60))

		}
	}
}
