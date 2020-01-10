package backend

import (
	"math/rand"
	"time"

	"github.com/go-qamel/qamel"
)

// BackEnd is the bridge for communicating between QML and Go
type BackEnd struct {
	qamel.QmlObject
	_ func() int   `slot:"getRandomNumber"`
	_ func(string) `signal:"timeChanged"`
	_ func()       `slot:"startTimer"`
}

func (b *BackEnd) getRandomNumber() int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(9999)
}

func (b *BackEnd) startTimer() {
	go func() {
		for {
			now := time.Now().Format("15:04:05")
			b.timeChanged(now)
			time.Sleep(time.Second)
		}
	}()
}
