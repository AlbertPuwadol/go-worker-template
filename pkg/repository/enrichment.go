package repository

import (
	"github.com/AlbertPuwadol/go-worker-template/pkg/adapter"
	"github.com/AlbertPuwadol/go-worker-template/pkg/entity"
	formatter "github.com/wisesight/spider-go-formatter"
)

type Enrichment interface {
	GetEnrichment(method entity.Method, post *formatter.Spider) error
}

type enrichement struct {
	enrichmentAdapter adapter.GRPC
}

func NewEnrichment(enrichmentAdapter adapter.GRPC) *enrichement {
	return &enrichement{enrichmentAdapter: enrichmentAdapter}
}

func (e enrichement) GetEnrichment(method entity.Method, post *formatter.Spider) error {
	postWithMutex := entity.PostWithMutex{Post: post}

	switch method {
	case entity.NER:
		nerResult, err := e.enrichmentAdapter.GetNer()
		if err != nil {
			return err
		}
		postWithMutex.SetNER(nerResult)
	case entity.Sentiment:
		sentimentResult, err := e.enrichmentAdapter.GetSentiment()
		if err != nil {
			return err
		}
		postWithMutex.SetSentiment(sentimentResult)
	case entity.TFIDF:
		tfidfResult, err := e.enrichmentAdapter.GetTfIdf()
		if err != nil {
			return err
		}
		postWithMutex.SetTfIdf(tfidfResult)
	}

	return nil
}
