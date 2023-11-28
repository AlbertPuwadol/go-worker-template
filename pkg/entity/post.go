package entity

import (
	"sync"

	formatter "github.com/wisesight/spider-go-formatter"
)

type PostWithMutex struct {
	mutex sync.Mutex
	Post  *formatter.Spider
}

func (p *PostWithMutex) SetNER(nerResult map[string]interface{}) {
	p.mutex.Lock()
	p.Post.KirinX.Ner = nerResult
	var temp []interface{}
	temp = make([]interface{}, 0)
	for _, v := range nerResult["entities"].([]map[string]interface{}) {
		temp = append(temp, v)
	}
	p.Post.KirinResult.Ner = temp
	p.mutex.Unlock()
}

func (p *PostWithMutex) SetSentiment(sentimentResult map[string]interface{}) {
	p.mutex.Lock()
	p.Post.KirinResult.Ocr += sentimentResult["sentiment"].(string)
	p.mutex.Unlock()
}

func (p *PostWithMutex) SetTfIdf(tfidfResult map[string]interface{}) {
	p.mutex.Lock()
	p.Post.KirinResult.Ocr += tfidfResult["tfidf"].(string)
	p.mutex.Unlock()
}
