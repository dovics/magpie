package scheduler

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

type Worker struct {
	name     string
	reporter Reporter
	taskCh   chan Task
	stopCh   chan struct{}
}

func (w *Worker) Run() {
	for {
		select {
		case task := <-w.taskCh:
			ctx := context.Background()
			result, err := task.Do(ctx)
			if err != nil {
				w.reporter.Fail(ctx, task.Name(), []byte(err.Error()))
				continue
			}

			if err := w.reporter.Success(ctx, task.Name(), result); err != nil {
				log.Printf("reporter success failed: %s\n", err)
				continue
			}
		case <-w.stopCh:
			return
		}

	}
}

func (w *Worker) Stop() {
	close(w.stopCh)
	close(w.taskCh)
}

func (w *Worker) RecvTask(ctx context.Context, task Task) error {
	w.taskCh <- task
	return nil
}

type Task interface {
	Name() string
	Resource() *Resource
	Do(context.Context) ([]byte, error)
}

type containerTask struct {
	name             string
	config           *container.Config
	hostConfig       *container.HostConfig
	networkingConfig *network.NetworkingConfig
}

func (t *containerTask) Do(ctx context.Context) ([]byte, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	reader, err := cli.ImagePull(ctx, t.config.Image, types.ImagePullOptions{})
	if err != nil {
		return nil, err
	}
	io.Copy(os.Stdout, reader)

	resp, err := cli.ContainerCreate(ctx, t.config, t.hostConfig, t.networkingConfig, nil, t.name)
	if err != nil {
		return nil, err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return nil, err
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return nil, err
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(out)
}
