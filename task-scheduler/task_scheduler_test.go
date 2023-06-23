package main

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestScheduler_Run(t *testing.T) {

	tasks := []*Task{
		{
			taskID:    "id1",
			taskName:  "task-1",
			taskFunc:  func() error { return nil },
			startTime: time.Now().Add(-1 * time.Second),
			status:    SCHEDULED,
		}, {
			taskID:    "id2",
			taskName:  "task-2",
			taskFunc:  func() error { return nil },
			startTime: time.Now().Add(10 * time.Second),
			status:    SCHEDULED,
		}, {
			taskID:    "id3",
			taskName:  "task-3",
			taskFunc:  func() error { return nil },
			startTime: time.Now().Add(1 * time.Second),
			status:    SCHEDULED,
		},
	}
	taskMap := make(map[string]*Task)
	for _, task := range tasks {
		taskMap[task.taskID] = task
	}

	scheduler := Scheduler{
		lock:    sync.Mutex{},
		done:    make(chan os.Signal),
		taskMap: taskMap,
	}

	go scheduler.Run()

	time.Sleep(time.Second * 2)
	scheduledCount := 0
	for _, task := range scheduler.taskMap {
		if task.status == SCHEDULED {
			scheduledCount++
		}
	}
	assert.Equal(t, 1, scheduledCount)
	if scheduledCount > 1 {
		t.Errorf("expected count is 1  but got %d", scheduledCount)
	}

	time.Sleep(time.Second * 10)

	for _, task := range scheduler.taskMap {
		if task.status == SCHEDULED {
			t.Errorf("Expected to be compelted")
		}
	}
}
