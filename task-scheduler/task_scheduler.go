package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

/**


Start time : 2:35 PM
Design task scheduler
Atlassian wants to build a simple scheduling system that can execute one-off HTTP requests at a specified time.
Consumers are internal Atlassian services only, and they may use this system to schedule requests to be executed in the future.
For example: A consumer can schedule a request to https://testservice.com/api/items/_create to be executed on upcoming
Sunday at 1am. So on upcoming Sunday at 1am the scheduling system must execute the HTTP request.
Requirements:
 Provide an interface for scheduling HTTP requests at a specified time
		- Execute an HTTP request at specified time.


Create task
	- taskID
	- taskName
	- task func()
	- execution_start_time
	- status

Scheduler
	- []tasks
	startScheduler()
			- scans the tasks and check if any task needs to be executed in the current time, if yes execute.

	How do we capture results??
		- For now, assume we are printing the results to system console.




Extension :
		- have a way to print the result to a file with the task name
		- Instead of scanning all the records can we optimize the scanning of the tasks?


*/

const (
	COMPLETED  = "completed"
	SCHEDULED  = "scheduled"
	INPROGRESS = "in-progress"
	FAILED     = "failed"
)

type Task struct {
	taskID    string
	taskName  string
	taskFunc  func() error
	startTime time.Time
	status    string
}

type Scheduler struct {
	lock  sync.Mutex
	taskMap map[string]*Task
	done chan os.Signal
}

func (s *Scheduler) AddTask(t *Task) {
	s.lock.Lock()
	s.taskMap[t.taskID] = t
	s.lock.Unlock()
}

func(s *Scheduler) GetTotalTasks() int{
	return len(s.taskMap)
}

func (s *Scheduler) RemoveTask(taskId string) {
	s.lock.Lock()
	delete(s.taskMap, taskId)
	s.lock.Unlock()
}

func (s *Scheduler) Run() {
	fmt.Println("Starting the scheduler.... ")
	signal.Notify(s.done, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)
	ticker := time.NewTicker(time.Second)
	stopScheduler := false
	wg := sync.WaitGroup{}
	for !stopScheduler {
		select {
		case <-ticker.C:
			s.lock.Lock()
			for _, task := range s.taskMap {
				if task.status == SCHEDULED && task.startTime.Before(time.Now()){
					task.status = INPROGRESS
					wg.Add(1)
					go executeTask(task, &wg)
				}
			}
			s.lock.Unlock()
		case <-s.done:
			fmt.Println("Force Stopped ")
			stopScheduler = true
			break

		}
	}
	wg.Wait()
	fmt.Println("Shutting down after gracefully completing all tasks")

}

func executeTask(task *Task, wg *sync.WaitGroup) {
	time.Sleep(time.Second * 1)
	err := task.taskFunc()
	if err != nil {
		fmt.Println(err)
		task.status = FAILED
	}else{
		task.status = COMPLETED
	}
	wg.Done()


}

func main(){

	taskArr := make([]*Task,0)
	for i := 0; i < 5; i++ {
		task := Task{
			taskID:   strconv.Itoa(i) ,
			taskName:  fmt.Sprintf("task-%d", i),
			taskFunc: func(i int) func() error {
				return func() error {
					fmt.Printf("This is task-%d\n", i)
					return nil
				}
			}(i),
			startTime: time.Now().Add(time.Second * time.Duration(i)),
			status:    SCHEDULED,
		}
		taskArr = append(taskArr, &task)
	}
	scheduler := Scheduler{
		lock:  sync.Mutex{},
		taskMap: make(map[string]*Task),
		done: make(chan os.Signal),
	}
	scheduler.Run()


}