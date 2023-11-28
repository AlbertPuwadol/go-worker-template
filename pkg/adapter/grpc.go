package adapter

import (
	"time"
)

type GRPC interface {
	GetNer() (map[string]interface{}, error)
	GetSentiment() (map[string]interface{}, error)
	GetTfIdf() (map[string]interface{}, error)
}

type gRPC struct {
}

func NewGRPC() *gRPC {
	return &gRPC{}
}

func (g gRPC) GetNer() (map[string]interface{}, error) {
	time.Sleep(5 * time.Second)
	var res map[string]interface{}
	res = make(map[string]interface{})
	res["preprocessed_text"] = "test"
	res["entities"] = []map[string]interface{}{{"label": "ORG", "entities": []string{"สุขสนุกเซ็นเตอร์", "some org"}}, {
		"label":    "LOC",
		"entiites": []string{"เชียงใหม่"},
	}}
	return res, nil
}

func (g gRPC) GetSentiment() (map[string]interface{}, error) {
	time.Sleep(1 * time.Second)
	var res map[string]interface{}
	res = make(map[string]interface{})
	res["preprocessed_text"] = "test"
	res["sentiment"] = "Positive"
	return res, nil
}

func (g gRPC) GetTfIdf() (map[string]interface{}, error) {
	time.Sleep(10 * time.Second)
	var res map[string]interface{}
	res = make(map[string]interface{})
	res["preprocessed_text"] = "test"
	res["tfidf"] = "0.99"
	return res, nil
}
