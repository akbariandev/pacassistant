package scheduler

import (
	"context"
	"sync"

	"github.com/akbariandev/pacassistant/internal/job"
)

type Scheduler interface {
	Submit(job job.Job)
	Run()
	Shutdown()
}

type scheduler struct {
	jobs   []job.Job
	mu     *sync.Mutex
	wg     *sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

func NewScheduler() Scheduler {
	return &scheduler{
		jobs: make([]job.Job, 0),
		mu:   &sync.Mutex{},
		wg:   &sync.WaitGroup{},
	}
}

func (s *scheduler) Submit(job job.Job) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.jobs = append(s.jobs, job)
}

func (s *scheduler) Run() {
	for _, j := range s.jobs {
		s.wg.Add(1)
		go j.Start()
	}
}

func (s *scheduler) Shutdown() {
	for _, j := range s.jobs {
		j.Stop()
		s.wg.Done()
	}
}
