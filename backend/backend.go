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
var today time.Time
var displayed_date time.Time

func init() {
	time_layout = "03:04PM"
	// os.Remove("/home/steams/Development/tempus/tempus_cli.db")

	db, err := sqlx.Open("sqlite3", "/home/steams/Development/tempus/tempus.db")
	// db, err := sqlx.Open("sqlite3", "/home/steams/Development/tempus/tempus_cli.db")

	if err != nil {
		panic(err)
	}

	// db.MustExec(activity_repo.Schema)
	repo := activity_repo.New(db)

	service = activity.CreateService(repo)

	today = time.Now()
	displayed_date = time.Now()
}

type Backend struct {
	qamel.QmlObject
	_ func(string)                                           `signal:"timeChanged"`
	_ func()                                                 `signal:"signalPause"`
	_ func()                                                 `signal:"signalStop"`
	_ func()                                                 `signal:"signalStart"`
	_ func(string, string, string, string, float64)          `signal:"updateList"`
	_ func(string, string, float64, string, float64, string) `signal:"updateTimeline"`
	_ func(string, float64)                                  `signal:"updateReport"`
	_ func(string)                                           `signal:"tagAdded"`
	_ func()                                                 `signal:"clearList"`
	_ func()                                                 `signal:"clearTimeline"`
	_ func()                                                 `signal:"clearReports"`
	_ func(string)                                           `signal:"dateChanged"`
	_ func(string, string)                                   `slot:"toggleStart"`
	_ func()                                                 `slot:"changeActivity"`
	_ func(string)                                           `slot:"changeTask"`
	_ func()                                                 `slot:"load"`
	_ func()                                                 `slot:"dateBack"`
	_ func()                                                 `slot:"dateForward"`
	_ func(string)                                           `slot:"addTag"`

	is_running   bool
	is_paused    bool
	timer        activity.Timer
	stopper      chan int
	current_tags []string
}

func duration(tasks []activity.Task) float64 {
	var sum time.Duration = 0
	for _, x := range tasks {
		sum = sum + x.End.Sub(x.Start)
	}
	return sum.Seconds()
}

func (b *Backend) dispatchListUpdate() {

	b.clearList()

	stuff := service.GetTasksByDay(displayed_date)
	// fmt.Println(stuff)
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

	sessions := service.GetTasksByDay(displayed_date)
	activities := make(map[string][]activity.Task)

	calc_duration := func(tasks []activity.Task) time.Duration {

		var dur time.Duration = 0

		for _, x := range tasks {
			dur = dur + x.End.Sub(x.Start)
		}

		return dur
	}

	for _, x := range sessions {
		if val, ok := activities[x.Act_name]; ok {
			activities[x.Act_name] = append(val, x.Tasks...)
		} else {
			activities[x.Act_name] = x.Tasks
		}
	}

	// fmt.Println("-----------------------Activities")
	// fmt.Println(activities)
	// fmt.Println("-----------------------")

	b.clearReports()
	for k, v := range activities {
		b.updateReport(
			k,
			calc_duration(v).Hours(),
		)
	}
}

func (b *Backend) dispatchTimelineUpdate() {
	b.clearTimeline()

	stuff := service.GetTasksByDay(displayed_date)
	// fmt.Println("Timeline stuff")
	// fmt.Println(stuff)

	year, month, day := displayed_date.Date()
	day_start := time.Date(year, month, day, 0, 0, 0, 0, time.Now().Location())
	morning := day_start.Add(time.Hour * 8)

	// fmt.Println(morning.Format(time_layout))

	for _, i := range stuff {
		for _, x := range i.Tasks {

			b.updateTimeline(
				x.Start.Format(time_layout),
				x.End.Format(time_layout),
				x.End.Sub(x.Start).Hours(),
				i.Act_name+" : "+x.Name,
				x.Start.Sub(morning).Hours(),
				i.Act_name,
			)
		}
	}

}
func (b *Backend) addTag(tag string) {
	b.current_tags = append(b.current_tags, tag)
	if b.timer != nil {
		b.timer.SetTags(b.current_tags)
	}
	b.tagAdded(tag)
}

func (b *Backend) dateForward() {
	displayed_date = displayed_date.AddDate(0, 0, 1)

	b.load()

	y1, m1, d1 := displayed_date.Date()

	y2, m2, d2 := time.Now().Date()

	if y1 == y2 && m1 == m2 && d1 == d2 {
		b.dateChanged("Today")
		return
	}

	y2, m2, d2 = time.Now().AddDate(0, 0, -1).Date()

	if y1 == y2 && m1 == m2 && d1 == d2 {
		b.dateChanged("Yesterday")
		return
	}

	b.dateChanged(displayed_date.Format("Mon, 02 Jan"))
	fmt.Println(time.UnixDate)
}

func (b *Backend) dateBack() {
	displayed_date = displayed_date.AddDate(0, 0, -1)

	b.load()
	y1, m1, d1 := displayed_date.Date()

	y2, m2, d2 := time.Now().Date()

	if y1 == y2 && m1 == m2 && d1 == d2 {
		b.dateChanged("Today")
		return
	}

	y2, m2, d2 = time.Now().AddDate(0, 0, -1).Date()

	if y1 == y2 && m1 == m2 && d1 == d2 {
		b.dateChanged("Yesterday")
		return
	}

	b.dateChanged(displayed_date.Format("Mon, 02 Jan"))
}

func (b *Backend) load() {
	b.dispatchListUpdate()
	b.dispatchReportUpdate()
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
	b.pause()
	b.timer.NewTask(name)

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
	b.timer.SetTags(b.current_tags)
	b.is_running = true
	b.stopper = make(chan int)

	// fmt.Println(b.timer)

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
			// fmt.Println("Stopping")
			return
		default:
			time.Sleep(time.Second)
			b.timer.UpdateDuration()
			// fmt.Println(b.timer)
			t := b.timer.GetDuration()
			b.timeChanged(fmt.Sprintf("%d h %d m %d s", int(t.Hours()), int(t.Minutes())%60, int(t.Seconds())%60))

		}
	}
}
