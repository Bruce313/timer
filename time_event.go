package main

import (
    "time"
)

//TimeEvent contain key and data of time
type TimeEvent struct {
    key string
    data []byte
    time time.Duration
}
