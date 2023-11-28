package usecase

import (
	"sync"

	"github.com/AlbertPuwadol/go-worker-template/pkg/entity"
	"github.com/AlbertPuwadol/go-worker-template/pkg/repository"
	formatter "github.com/wisesight/spider-go-formatter"
)

type Enrichment interface {
	GetPostEnrichment(post formatter.Spider) error
}

type enrichement struct {
	enrichmentRepository repository.Enrichment
}

func NewEnrichment(enrichmentRepository repository.Enrichment) *enrichement {
	return &enrichement{enrichmentRepository: enrichmentRepository}
}

func (e enrichement) GetPostEnrichment(post *formatter.Spider) error {
	done := make(chan struct{})
	defer close(done)

	var enrichmentMethodList = []entity.Method{entity.NER, entity.Sentiment, entity.TFIDF}
	errorChannel := make(chan error, len(enrichmentMethodList))

	var enrichmentMethodListChannel = make(chan entity.EnrichementMethod)

	go func() {
		defer close(enrichmentMethodListChannel)
		for _, enrichmentMethod := range enrichmentMethodList {
			select {
			case enrichmentMethodListChannel <- entity.EnrichementMethod{Method: enrichmentMethod, Post: post}:
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
			for enrichmentMethod := range enrichmentMethodListChannel {
				select {
				case errorChannel <- e.enrichmentRepository.GetEnrichment(enrichmentMethod.Method, enrichmentMethod.Post):
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
