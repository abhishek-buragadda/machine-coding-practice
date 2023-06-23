package main

import (
	"fmt"
	"sort"
	"time"
)

/*
	requirements:
		= implement a event scheduler .
	What does that do :
		- We have events .
		- We have a scheduler whose job is to execute the event when the start time is reached.
		- The scheduler should check timings of all the events and execute them when the time comes.

*/

type IEvent interface {
	Action() error
}

type Event struct {
	eventID   string
	eventName string
	startTime time.Time
	endTime   time.Time
	completed bool
}

func (event *Event) Action() error {
	fmt.Printf("Performing action for %s", event.eventID)
	event.completed = true
	return nil
}

type EventScheduler struct {
	ID     string
	Events []Event
	Done   chan struct{}
}

func (es *EventScheduler) AddEvent(e Event) {
	es.Events = append(es.Events, e)
	sort.Slice(es.Events, func(i, j int) bool {
		return es.Events[i].startTime.Before(es.Events[j].startTime)
	})
}

func (es *EventScheduler) Run() {

	ticker := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-ticker.C:

			for index, event := range es.Events {
				if event.startTime.Before(time.Now()) && event.endTime.After(time.Now()) && !event.completed {
					err := event.Action()
					if err == nil {
						es.Events[index].completed = true
					}
				}
			}
			count := 0
			for _, event := range es.Events {
				if !event.completed {
					count++
				}
			}
			if count == 0 {
				ticker.Stop()
				fmt.Println("Ticker stopped")
				es.Done <- struct{}{}
			}
		}
	}

}

//func main() {
//	e1 := Event{
//		eventID:   "1",
//		eventName: "one",
//		startTime: time.Now().Add(-time.Minute),
//		endTime:   time.Now().Add(2 * time.Minute),
//	}
//	e2 := Event{
//		eventID:   "2",
//		eventName: "two",
//		startTime: time.Now().Add(5 * time.Second),
//		endTime:   time.Now().Add(time.Minute),
//	}
//	scheduler := EventScheduler{
//		ID:     "schedule-1",
//		Events: []Event{},
//		Done:   make(chan struct{}),
//	}
//
//	scheduler.AddEvent(e1)
//	scheduler.AddEvent(e2)
//
//	go scheduler.Run()
//	<-scheduler.Done
//	fmt.Println("Scheduler stopped")
//
//}
