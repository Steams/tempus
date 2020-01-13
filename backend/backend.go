package backend

import (
	"fmt"
	"tempus/pkg/activity"
	"time"

	"github.com/go-qamel/qamel"
)

// func init() {
// 	state = model{false, time.Now(), time.Now(), "", []string{}, make(chan int)}
// }

// BackEnd is the bridge for communicating between QML and Go
type BackEnd struct {
	qamel.QmlObject
	_ func(string)         `signal:"timeChanged"`
	_ func(string, string) `slot:"start"`
	_ func()               `slot:"pause"`
	_ func() bool          `slot:"isRunning"`

	is_running bool
	timer      activity.Timer
	stopper    chan int
}

// func (b *BackEnd) isRunning() bool {
// 	return false
// }
func (b *BackEnd) isRunning() bool {
	return b.is_running
}

func (b *BackEnd) start(activity_name, task string) {

	b.timer = activity.NewTimer(activity_name, task)
	b.is_running = true
	b.stopper = make(chan int)

	// // a := activity.ActivitySpan{task, time.Now(), time.Now(), nil}
	fmt.Println(b.timer)

	go func() {
		for {
			select {
			case <-b.stopper:
				b.is_running = false
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
	}()
}

func (b *BackEnd) pause() {
	// a := activity.ActivitySpan{task, time.Now(), time.Now(), nil}
	fmt.Println("Pausing")

	b.timer.Pause()
	b.stopper <- 1
}
