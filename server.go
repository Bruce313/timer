package main

import (
    "time"
    "fmt"
)

func main() {
    chTimeEventOut := make(chan *TimeEvent)
    chTimeEventCmd := make(chan *TimeEventCmd)
    trigger := NewTimeEventTrigger(chTimeEventOut, chTimeEventCmd)
    go trigger.Begin()
    sec := &TimeEvent {
        key: "a",
        data: []byte("a"),
        delay: time.Second * 2,
    }
    trigger.addTimeEvent(sec)
    te := <-chTimeEventOut
    fmt.Printf("%s\n", te)
}
