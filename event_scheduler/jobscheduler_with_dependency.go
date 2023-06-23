package main

import (
	"fmt"
	"strconv"
	"sync"
)

/**

		4:51 PM

		Implement a job scheduler with dependency;

		How is the dependency maintained?
			- for each job we have list of dependencies wihch needs to be executeed before the job.
		Assumption:
			 Job  -> jobId, dependency,  task

			scheduler

				[]jobs
				Run()
					- execute the jobs in parallel.
						- for each job, execute the list of dependencies and once that is done, execute the task at the end.
					- Wait for the jobs to complete.
						- wait for the run() method until all the jobs are executed.



 */

type Job struct {
	jobId int
	dependency []*Job
	task func(jobId int)
}

type Scheduler struct {
	jobs map[int]*Job
	schedulerID string

}

func(s *Scheduler) AddJob(job *Job){
	s.jobs[job.jobId] = job
}

func(s *Scheduler) AddJobWithDependency(jobId int, dependencyId int){

	currJob :=  s.jobs[jobId]
	if dependentJob, ok := s.jobs[dependencyId]; ok {
		currJob.dependency = append(currJob.dependency, dependentJob )
	}


	delete(s.jobs, dependencyId)
}


func (s *Scheduler) Run(n int) {
	wg := sync.WaitGroup{}
	wg.Add(n)


	for _, job := range s.jobs {
			go func(j *Job){
				s.runJob(j, &wg)
			}(job)
	}
	wg.Wait()
}

func (s *Scheduler) runJob(j *Job, wg *sync.WaitGroup) {

	if len(j.dependency) != 0 {
		for i := 0; i < len(j.dependency); i++ {
			s.runJob(j.dependency[i], wg)
		}

	}
	j.task(j.jobId)
	wg.Done()
}





func main() {

	n := 8
	jobMap := make(map[int]*Job)
	for i := 0; i < n; i++ {
		j := &Job{
			jobId:     i,
			dependency: []*Job{},
			task: func(k int) {
				fmt.Printf("task for job :%s \n", strconv.Itoa(k))
			},
		}
		jobMap[i] = j

	}


	s := Scheduler{
		jobs:        jobMap,
		schedulerID: "scheduler-1",
	}
	s.AddJobWithDependency(1, 5)
	s.AddJobWithDependency(2, 6)
	s.AddJobWithDependency(3, 7)
	s.Run(n)
	fmt.Printf("completed execution of all the jobs")

}