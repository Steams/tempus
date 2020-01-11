package activity

import "time"

type ActivitySpan struct {
	Name  string
	Start time.Time
	End   time.Time
	Tasks []TaskSpan
}

type TaskSpan struct {
	Name  string
	Start string
	End   string
}

type ActivityStarter struct {
	name  string
	start time.Time
}

// func StartTimer(activity,task string) {
// 	// var t time.Duration
// 	// t = 0 60 * time.Second

// 	a := activity.ActivitySpan{task, time.Now(), time.Now(), nil}

// 	fmt.Println(a)

// 	go func() {
// 		for {
// 			// now := time.Now().Format("15:04:05")
// 			// b.timeChanged(now)
// 			time.Sleep(time.Second)
// 			t = t + time.Second

// 			fmt.Println(int(t.Seconds()))
// 			b.timeChanged(fmt.Sprintf("%d hrs %d min %d s", int(t.Hours()), int(t.Minutes()), int(t.Seconds())))
// 		}
// 	}()
// }

// func PauseTimer() {

// }

// func UpdateTask() {

// }
