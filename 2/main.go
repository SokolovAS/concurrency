// Hint 1: time.Ticker can be used to cancel function
// Hint 2: to calculate time-diff for Advanced lvl use:
//  start := time.Now()
//	// your work
//	t := time.Now()
//	elapsed := t.Sub(start) // 1s or whatever time has passed

package main

import (
	"time"
)

var freeTime = time.Duration(10 * time.Second)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool          // can be used for 2nd level task. Premium users won't have 10 seconds limit.
	TimeUsed  time.Duration // in seconds
}

// HandleRequest runs the processes requested by users. Returns false if process had to be killed
func HandleRequest(process func(), u *User) bool {
	chStop := make(chan bool)
	chProc := make(chan bool)

	go func() {
		t := freeTime - u.TimeUsed
		stop := time.After(t)
		select {
		case <-stop:
			chStop <- true
		}
	}()

	go func() {
		start := time.Now()
		process()
		t := time.Now()
		elapsed := t.Sub(start)
		u.TimeUsed = elapsed
		chProc <- true
	}()

	select {
	case <-chStop:
		return false
	case <-chProc:
		return true
	}
}

func main() {
	RunMockServer()
}
