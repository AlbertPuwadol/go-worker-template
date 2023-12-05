package entity

type Method string

const (
	Task1 Method = "task1"
	Task2 Method = "task2"
	Task3 Method = "task3"
)

type TaskMethod struct {
	Method Method
	Post   *Post
}
