//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer szenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"fmt"
	"time"
)

var tweetChannel = make(chan Tweet)
var doneChannel = make(chan struct{})

func producer(stream Stream) {
	for {
		tweet, err := stream.Next()
		if err == ErrEOF {
			close(tweetChannel)
			return
		}
		tweetChannel <- *tweet
	}
}

func consumer() {
	for {
		t, ok := <-tweetChannel
		if ok {
			if t.IsTalkingAboutGo() {
				fmt.Println(t.Username, "\ttweets about golang")
			} else {
				fmt.Println(t.Username, "\tdoes not tweet about golang")
			}
		} else {
			doneChannel <- struct{}{}
			return
		}
	}

}

func main() {
	start := time.Now()
	stream := GetMockStream()

	// Producer
	go producer(stream)

	// Consumer
	go consumer()

	<-doneChannel

	fmt.Printf("Process took %s\n", time.Since(start))
}
