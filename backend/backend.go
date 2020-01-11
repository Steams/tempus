package backend

import (
	"fmt"
	"time"

	"github.com/go-qamel/qamel"
)

var state model

func init() {
	state = model{false, time.Now(), time.Now(), "", []string{}, make(chan int)}
}

type model struct {
	is_running bool
	start      time.Time
	end        time.Time
	activity   string
	tasks      []string
	stopper    chan int
}

// BackEnd is the bridge for communicating between QML and Go
type BackEnd struct {
	qamel.QmlObject
	_ func(string)         `signal:"timeChanged"`
	_ func(string, string) `slot:"startTimer"`
	_ func()               `slot:"pauseTimer"`
	_ func() bool          `slot:"isRunning"`
}

func (b *BackEnd) isRunning() bool {
	return state.is_running
}

func (b *BackEnd) startTimer(activity, task string) {

	state.activity = activity
	state.tasks = append(state.tasks, task)
	state.start = time.Now()
	state.is_running = true

	// a := activity.ActivitySpan{task, time.Now(), time.Now(), nil}
	fmt.Println(state)

	go func() {
		var t time.Duration = time.Second * 0
		for {
			select {
			case <-state.stopper:
				state.is_running = false
				fmt.Println("Stopping")
				return
			default:
				// now := time.Now().Format("15:04:05")
				// b.timeChanged(now)
				time.Sleep(time.Second)
				t = t + time.Second
				fmt.Println(int(t.Seconds()))
				b.timeChanged(fmt.Sprintf("%d hrs %d min %d s", int(t.Hours()), int(t.Minutes()), int(t.Seconds())))

			}
		}
	}()
}

func (b *BackEnd) pauseTimer() {
	// a := activity.ActivitySpan{task, time.Now(), time.Now(), nil}
	state.stopper <- 1

	fmt.Println(state)
}
