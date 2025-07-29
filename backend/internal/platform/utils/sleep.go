package utils

import (
	"math/rand"
	"time"
)

// RandomSleep sleeps for a random duration between 1 and 2 seconds.
func RandomSleep() {
	rand.Seed(time.Now().UnixNano())
	sleepDuration := time.Duration(100+rand.Intn(100)) * time.Millisecond // 1s to 2s
	time.Sleep(sleepDuration)
}
