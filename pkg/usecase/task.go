package usecase

import (
	"sync"

	"github.com/AlbertPuwadol/go-worker-template/pkg/entity"
	"github.com/AlbertPuwadol/go-worker-template/pkg/repository"
)

type Task interface {
	GetPostTasks(post *entity.Post) error
}

type task struct {
	taskRepository repository.Task
}

func NewTask(taskRepository repository.Task) *task {
	return &task{taskRepository: taskRepository}
}

func (e task) GetPostTasks(post *entity.Post) error {
	done := make(chan struct{})
	defer close(done)

	var taskMethodList = []entity.Method{entity.Task1, entity.Task2, entity.Task3}
	errorChannel := make(chan error, len(taskMethodList))

	var taskMethodListChannel = make(chan entity.TaskMethod)

	go func() {
		defer close(taskMethodListChannel)
		for _, taskMethod := range taskMethodList {
			select {
			case taskMethodListChannel <- entity.TaskMethod{Method: taskMethod, Post: post}:
			case <-done:
			}
		}
	}()

	var wg sync.WaitGroup
	const CONCURRENT = 2

	wg.Add(CONCURRENT)
	for i := 0; i < CONCURRENT; i++ {
		go func() {
			defer wg.Done()
			for taskMethod := range taskMethodListChannel {
				select {
				case errorChannel <- e.taskRepository.GetTasks(taskMethod.Method, taskMethod.Post):
				case <-done:
				}
			}
		}()
	}

	wg.Wait()
	close(errorChannel)
	for err := range errorChannel {
		if err != nil {
			return err
		}
	}

	return nil
}
