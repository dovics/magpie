package scheduler

import (
	"context"
)

type DiscoverManager interface {
	RecvCh() chan []*WorkerConfig
}

type Scheduler struct {
	discoverManager DiscoverManager
	dispatcher      Dispatcher

	storage Storage
	cache   *Cache
	taskCh  chan Task
	stopCh  chan struct{}
}

func (s *Scheduler) Run() {
	for {
		select {
		case workers := <-s.discoverManager.RecvCh():
			for _, worker := range workers {
				if _, ok := s.cache.nodes[worker.name]; ok {
					continue
				}

				s.cache.AddNode(worker)
			}
		case task := <-s.taskCh:
			for _, n := range s.cache.nodes {
				s.cache.AddTask(n.name, task)
				ctx := context.TODO()
				s.dispatcher.Dispatch(ctx, n.name, task)
				break
			}
		}
	}
}

func (s *Scheduler) Stop() {
	close(s.stopCh)
	close(s.taskCh)
}

func (s *Scheduler) AddTask(task Task) error {
	s.taskCh <- task
	return nil
}

func (s *Scheduler) Start(ctx context.Context, name string) error {
	s.cache.UpdateTaskStatus(name, Running)
	return nil
}

func (s *Scheduler) Success(ctx context.Context, name string, result []byte) error {
	s.cache.UpdateTaskStatus(name, Completed)
	return s.storage.Save(name, result)
}

func (s *Scheduler) Fail(ctx context.Context, name string, reason []byte) error {
	s.cache.UpdateTaskStatus(name, Failed)
	return s.storage.Save(name, reason)
}
