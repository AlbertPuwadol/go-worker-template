package repository

import (
	"github.com/AlbertPuwadol/go-worker-template/pkg/adapter"
	"github.com/AlbertPuwadol/go-worker-template/pkg/entity"
)

type Task interface {
	GetTasks(method entity.Method, post *entity.Post) error
}

type task struct {
	grpcAdapter adapter.GRPC
}

func NewTask(grpcAdapter adapter.GRPC) *task {
	return &task{grpcAdapter: grpcAdapter}
}

func (t task) GetTasks(method entity.Method, post *entity.Post) error {
	postWithMutex := entity.PostWithMutex{Post: post}

	switch method {
	case entity.Task1:
		task1Result, err := t.grpcAdapter.GetTask1(post.Text)
		if err != nil {
			return err
		}
		postWithMutex.SetTask1(task1Result)
	case entity.Task2:
		task2Result, err := t.grpcAdapter.GetTask2(post.Text)
		if err != nil {
			return err
		}
		postWithMutex.SetTask2(task2Result.Task2)
	case entity.Task3:
		task3Result, err := t.grpcAdapter.GetTask3(post.Text)
		if err != nil {
			return err
		}
		postWithMutex.SetTask3(int(task3Result.Task3))
	}

	return nil
}
