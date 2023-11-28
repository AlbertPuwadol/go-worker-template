package adapter

import (
	"time"
)

type KiringRPC interface {
	GetNer() (map[string]interface{}, error)
	GetSentiment() (map[string]interface{}, error)
	GetTfIdf() (map[string]interface{}, error)
}

type kiringRPC struct {
}

func NewKiringRPC() *kiringRPC {
	return &kiringRPC{}
}

func (k kiringRPC) GetNer() (map[string]interface{}, error) {
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

func (k kiringRPC) GetSentiment() (map[string]interface{}, error) {
	time.Sleep(1 * time.Second)
	var res map[string]interface{}
	res = make(map[string]interface{})
	res["preprocessed_text"] = "test"
	res["sentiment"] = "Positive"
	return res, nil
}

func (k kiringRPC) GetTfIdf() (map[string]interface{}, error) {
	time.Sleep(10 * time.Second)
	var res map[string]interface{}
	res = make(map[string]interface{})
	res["preprocessed_text"] = "test"
	res["tfidf"] = "0.99"
	return res, nil
}
