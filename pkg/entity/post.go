package entity

import (
	"sync"

	pb "github.com/AlbertPuwadol/grpc-clean/proto"
)

type Post struct {
	ID    string      `json:"id"`
	Text  string      `json:"text"`
	Task1 interface{} `json:"task1"`
	Task2 string      `json:"task2"`
	Task3 int         `json:"task3"`
}

type PostWithMutex struct {
	mutex sync.Mutex
	Post  *Post
}

func (p *PostWithMutex) SetTask1(task1Result *pb.Task1Response) {
	p.mutex.Lock()
	p.Post.Task1 = task1Result
	p.mutex.Unlock()
}

func (p *PostWithMutex) SetTask2(task2 string) {
	p.mutex.Lock()
	p.Post.Task2 = task2
	p.mutex.Unlock()
}

func (p *PostWithMutex) SetTask3(task3 int) {
	p.mutex.Lock()
	p.Post.Task3 = task3
	p.mutex.Unlock()
}
