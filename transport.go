package scheduler

import "context"

type Dispatcher interface {
	Dispatch(ctx context.Context, node string, task Task) error
}

type Reporter interface {
	Start(ctx context.Context, name string) error
	Success(ctx context.Context, name string, result []byte) error
	Fail(ctx context.Context, name string, reason []byte) error
}

type memTransport struct {
	scheduler *Scheduler
	workers   map[string]*Worker
}

func NewMemTransport(scheduler *Scheduler, workers []*Worker) *memTransport {
	workerMap := make(map[string]*Worker, len(workers))

	for _, worker := range workers {
		workerMap[worker.name] = worker
	}

	return &memTransport{
		scheduler,
		workerMap,
	}
}

func (t *memTransport) Dispatch(node string, task Task) error {
	return t.workers[node].RecvTask(context.TODO(), task)
}

func (t *memTransport) Start(ctx context.Context, name string) error {
	return t.scheduler.Start(ctx, name)
}

func (t *memTransport) Success(ctx context.Context, name string, result []byte) error {
	return t.scheduler.Success(ctx, name, result)
}

func (t *memTransport) Fail(ctx context.Context, name string, reason []byte) error {
	return t.scheduler.Fail(ctx, name, reason)
}
