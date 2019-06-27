//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"time"
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
}

func timeUp(c chan bool, u *User) {
	for {
		time.Sleep(time.Second)
		u.TimeUsed++

		if u.TimeUsed >= 10 {
			c <- false
		}
	}
	// time.Sleep(10 * time.Second)
	// c <- false
}

func processing(c chan bool, process func()) {
	process()
	c <- true
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	channel := make(chan bool)
	if !u.IsPremium {
		go timeUp(channel, u)
	}
	go processing(channel, process)

	return <-channel
}

func main() {
	RunMockServer()
}
